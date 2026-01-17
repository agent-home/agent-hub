package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	runLocal   bool
	runVersion string
	runInput   string
)

var runCmd = &cobra.Command{
	Use:   "run <namespace/name>",
	Short: "运行智能体",
	Long: `运行一个智能体进行交互。

示例:
  agenthub run agenthub/simple-assistant     # 交互模式
  agenthub run user/my-agent --local         # 运行本地智能体
  agenthub run user/my-agent -i "你好"        # 单次输入
  agenthub run user/my-agent@1.0.0           # 指定版本`,
	Args: cobra.ExactArgs(1),
	Run:  runAgent,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&runLocal, "local", "l", false, "运行本地智能体")
	runCmd.Flags().StringVarP(&runVersion, "version", "v", "latest", "指定版本")
	runCmd.Flags().StringVarP(&runInput, "input", "i", "", "直接输入 (非交互模式)")
}

func runAgent(cmd *cobra.Command, args []string) {
	agentRef := args[0]
	apiURL := viper.GetString("api_url")
	apiKey := viper.GetString("api_key")
	token := viper.GetString("token")

	namespace, name, version := parseAgentRef(agentRef)
	if version == "" {
		version = runVersion
	}

	if runLocal {
		runLocalAgent(namespace, name)
		return
	}

	fmt.Printf("🤖 启动 %s/%s@%s\n", namespace, name, version)
	fmt.Println("(输入 /exit 退出, /help 查看帮助)")
	fmt.Println()

	// 如果有直接输入，单次运行
	if runInput != "" {
		response := invokeAgent(apiURL, apiKey, token, namespace, name, runInput)
		fmt.Println(response)
		return
	}

	// 交互模式
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// 处理命令
		if strings.HasPrefix(input, "/") {
			switch input {
			case "/exit", "/quit", "/q":
				fmt.Println("再见！")
				return
			case "/help", "/h":
				printHelp()
				continue
			case "/clear":
				fmt.Print("\033[H\033[2J")
				continue
			default:
				fmt.Println("未知命令。输入 /help 查看帮助。")
				continue
			}
		}

		// 调用智能体
		fmt.Print("Agent: ")
		response := invokeAgent(apiURL, apiKey, token, namespace, name, input)
		fmt.Println(response)
		fmt.Println()
	}
}

func invokeAgent(apiURL, apiKey, token, namespace, name, input string) string {
	url := fmt.Sprintf("%s/invoke/%s/%s", apiURL, namespace, name)

	body, _ := json.Marshal(map[string]interface{}{
		"message": input,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	} else if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("调用失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("调用失败: HTTP %d", resp.StatusCode)
	}

	var result struct {
		Output string `json:"output"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Output
}

func runLocalAgent(namespace, name string) {
	home, _ := os.UserHomeDir()
	agentDir := filepath.Join(home, ".agenthub", "agents", namespace, name)

	// 查找最新版本
	versions, err := os.ReadDir(agentDir)
	if err != nil {
		fmt.Printf("未找到本地智能体: %s/%s\n", namespace, name)
		fmt.Println("使用 'agenthub pull' 先下载智能体")
		os.Exit(1)
	}

	if len(versions) == 0 {
		fmt.Println("未找到任何版本")
		os.Exit(1)
	}

	// 使用最后一个版本
	latestVersion := versions[len(versions)-1].Name()
	specPath := filepath.Join(agentDir, latestVersion, "agentspec.yaml")

	specData, err := os.ReadFile(specPath)
	if err != nil {
		fmt.Printf("读取 spec 失败: %v\n", err)
		os.Exit(1)
	}

	var spec struct {
		Runtime struct {
			Type string `yaml:"type"`
		} `yaml:"runtime"`
		Prompts struct {
			System string `yaml:"system"`
		} `yaml:"prompts"`
	}

	if err := yaml.Unmarshal(specData, &spec); err != nil {
		fmt.Printf("解析 spec 失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🤖 本地运行 %s/%s@%s\n", namespace, name, latestVersion)
	fmt.Printf("类型: %s\n", spec.Runtime.Type)
	fmt.Println()

	if spec.Runtime.Type == "prompt" {
		fmt.Println("System Prompt:")
		fmt.Println("---")
		fmt.Println(spec.Prompts.System)
		fmt.Println("---")
		fmt.Println()
		fmt.Println("这是一个纯提示词智能体，需要配合 LLM 使用。")
	} else {
		fmt.Println("本地运行需要相应的运行时环境。")
	}
}

func printHelp() {
	fmt.Println(`
可用命令:
  /exit, /quit, /q  退出
  /help, /h         显示帮助
  /clear            清屏
`)
}
