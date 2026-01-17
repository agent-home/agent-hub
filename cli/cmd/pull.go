package cmd

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pullVersion string
	pullOutput  string
)

var pullCmd = &cobra.Command{
	Use:   "pull <namespace/name>",
	Short: "下载智能体",
	Long: `从 AgentHub 下载智能体到本地。

示例:
  agenthub pull agenthub/code-reviewer
  agenthub pull user/my-agent@1.0.0
  agenthub pull user/my-agent -o ./my-agents/`,
	Args: cobra.ExactArgs(1),
	Run:  runPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().StringVarP(&pullVersion, "version", "v", "latest", "指定版本")
	pullCmd.Flags().StringVarP(&pullOutput, "output", "o", "", "输出目录")
}

func runPull(cmd *cobra.Command, args []string) {
	agentRef := args[0]
	apiURL := viper.GetString("api_url")

	// 解析 agent 引用
	namespace, name, version := parseAgentRef(agentRef)
	if version == "" {
		version = pullVersion
	}

	fmt.Printf("📥 正在下载 %s/%s@%s ...\n", namespace, name, version)

	// 获取版本信息
	versionURL := fmt.Sprintf("%s/api/v1/agents/%s/%s/versions/%s", apiURL, namespace, name, version)
	resp, err := http.Get(versionURL)
	if err != nil {
		fmt.Printf("下载失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("智能体 %s/%s@%s 不存在\n", namespace, name, version)
		os.Exit(1)
	}

	var versionInfo struct {
		Version string `json:"version"`
		Spec    string `json:"spec"`
		Size    int64  `json:"size"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&versionInfo); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		os.Exit(1)
	}

	// 确定输出目录
	outputDir := pullOutput
	if outputDir == "" {
		home, _ := os.UserHomeDir()
		outputDir = filepath.Join(home, ".agenthub", "agents", namespace, name, versionInfo.Version)
	}

	// 创建目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("创建目录失败: %v\n", err)
		os.Exit(1)
	}

	// 保存 spec 文件
	specPath := filepath.Join(outputDir, "agentspec.yaml")
	if err := os.WriteFile(specPath, []byte(versionInfo.Spec), 0644); err != nil {
		fmt.Printf("保存 spec 失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ 下载完成！\n")
	fmt.Printf("  位置: %s\n", outputDir)
	fmt.Printf("  版本: %s\n", versionInfo.Version)
	fmt.Printf("\n使用 'agenthub run %s/%s' 运行智能体\n", namespace, name)
}

// parseAgentRef 解析智能体引用
// 格式: namespace/name[@version]
func parseAgentRef(ref string) (namespace, name, version string) {
	// 检查版本
	if idx := strings.LastIndex(ref, "@"); idx != -1 {
		version = ref[idx+1:]
		ref = ref[:idx]
	}

	// 解析 namespace/name
	parts := strings.SplitN(ref, "/", 2)
	if len(parts) == 2 {
		namespace = parts[0]
		name = parts[1]
	} else {
		namespace = "agenthub" // 默认命名空间
		name = parts[0]
	}

	return
}

// extractTarGz 解压 tar.gz 文件
func extractTarGz(reader io.Reader, destDir string) error {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}
