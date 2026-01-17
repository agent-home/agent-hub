package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	pushVersion   string
	pushChangelog string
)

var pushCmd = &cobra.Command{
	Use:   "push [path]",
	Short: "发布智能体",
	Long: `将智能体发布到 AgentHub。

会读取目录中的 agentspec.yaml 文件并发布。

示例:
  agenthub push                        # 发布当前目录
  agenthub push ./my-agent             # 发布指定目录
  agenthub push -v 1.0.0               # 指定版本号
  agenthub push -m "修复了一些问题"      # 添加更新日志`,
	Run: runPush,
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringVarP(&pushVersion, "version", "v", "", "版本号 (覆盖 spec 中的版本)")
	pushCmd.Flags().StringVarP(&pushChangelog, "message", "m", "", "更新日志")
}

func runPush(cmd *cobra.Command, args []string) {
	// 检查登录状态
	token := viper.GetString("token")
	if token == "" {
		fmt.Println("请先登录: agenthub login")
		os.Exit(1)
	}

	// 确定路径
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// 读取 agentspec.yaml
	specPath := filepath.Join(path, "agentspec.yaml")
	specData, err := os.ReadFile(specPath)
	if err != nil {
		fmt.Printf("读取 agentspec.yaml 失败: %v\n", err)
		fmt.Println("确保当前目录包含 agentspec.yaml 文件")
		os.Exit(1)
	}

	// 解析 spec
	var spec struct {
		Version  string `yaml:"version"`
		Metadata struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Author      string `yaml:"author"`
			Category    string `yaml:"category"`
			Tags        []string `yaml:"tags"`
			License     string `yaml:"license"`
		} `yaml:"metadata"`
	}

	if err := yaml.Unmarshal(specData, &spec); err != nil {
		fmt.Printf("解析 agentspec.yaml 失败: %v\n", err)
		os.Exit(1)
	}

	// 确定版本号
	version := pushVersion
	if version == "" {
		version = spec.Version
	}
	if version == "" {
		version = "0.1.0"
	}

	username := viper.GetString("username")
	agentName := spec.Metadata.Name

	fmt.Printf("📤 正在发布 %s/%s@%s ...\n", username, agentName, version)

	apiURL := viper.GetString("api_url")

	// 1. 先检查或创建智能体
	checkURL := fmt.Sprintf("%s/api/v1/agents/%s/%s", apiURL, username, agentName)
	checkResp, err := http.Get(checkURL)
	if err != nil {
		fmt.Printf("检查智能体失败: %v\n", err)
		os.Exit(1)
	}

	if checkResp.StatusCode == http.StatusNotFound {
		// 创建智能体
		fmt.Println("  创建智能体...")
		createBody, _ := json.Marshal(map[string]interface{}{
			"name":        agentName,
			"description": spec.Metadata.Description,
			"category":    spec.Metadata.Category,
			"tags":        spec.Metadata.Tags,
			"license":     spec.Metadata.License,
			"visibility":  "public",
		})

		req, _ := http.NewRequest("POST", apiURL+"/api/v1/agents", bytes.NewBuffer(createBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		createResp, err := client.Do(req)
		if err != nil {
			fmt.Printf("创建智能体失败: %v\n", err)
			os.Exit(1)
		}

		if createResp.StatusCode != http.StatusCreated {
			fmt.Printf("创建智能体失败: HTTP %d\n", createResp.StatusCode)
			os.Exit(1)
		}
	}
	checkResp.Body.Close()

	// 2. 发布版本
	fmt.Println("  发布版本...")
	publishBody, _ := json.Marshal(map[string]interface{}{
		"version":   version,
		"spec":      string(specData),
		"changelog": pushChangelog,
	})

	publishURL := fmt.Sprintf("%s/api/v1/agents/%s/%s/versions", apiURL, username, agentName)
	req, _ := http.NewRequest("POST", publishURL, bytes.NewBuffer(publishBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	publishResp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发布版本失败: %v\n", err)
		os.Exit(1)
	}
	defer publishResp.Body.Close()

	if publishResp.StatusCode != http.StatusCreated {
		var errResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(publishResp.Body).Decode(&errResp)
		fmt.Printf("发布失败: %s\n", errResp.Error)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("✓ 发布成功！\n")
	fmt.Printf("  %s/%s@%s\n", username, agentName, version)
	fmt.Printf("\n查看: https://agenthub.dev/%s/%s\n", username, agentName)
}
