package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "初始化新智能体项目",
	Long: `在当前目录创建一个新的智能体项目。

会生成 agentspec.yaml 模板文件。

示例:
  agenthub init my-agent
  agenthub init                # 交互式创建`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	var name, description, category, runtimeType string

	if len(args) > 0 {
		name = args[0]
	} else {
		fmt.Print("智能体名称: ")
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)
	}

	if name == "" {
		fmt.Println("名称不能为空")
		os.Exit(1)
	}

	fmt.Print("描述: ")
	description, _ = reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("分类 (assistant/coding/writing/analysis/creative/education/business/research/tooling/other) [assistant]: ")
	category, _ = reader.ReadString('\n')
	category = strings.TrimSpace(category)
	if category == "" {
		category = "assistant"
	}

	fmt.Print("运行时类型 (prompt/python/nodejs/docker/remote) [prompt]: ")
	runtimeType, _ = reader.ReadString('\n')
	runtimeType = strings.TrimSpace(runtimeType)
	if runtimeType == "" {
		runtimeType = "prompt"
	}

	author := viper.GetString("username")
	if author == "" {
		author = "your-username"
	}

	// 生成 agentspec.yaml
	spec := generateSpec(name, description, category, runtimeType, author)

	// 写入文件
	specPath := "agentspec.yaml"
	if err := os.WriteFile(specPath, []byte(spec), 0644); err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("✓ 已创建 %s\n", specPath)
	fmt.Println()
	fmt.Println("下一步:")
	fmt.Println("  1. 编辑 agentspec.yaml 完善配置")
	if runtimeType == "prompt" {
		fmt.Println("  2. 编写系统提示词")
	} else if runtimeType == "python" {
		fmt.Println("  2. 创建 agent.py 和 requirements.txt")
	} else if runtimeType == "nodejs" {
		fmt.Println("  2. 创建入口文件和 package.json")
	}
	fmt.Println("  3. 运行 'agenthub push' 发布")
}

func generateSpec(name, description, category, runtimeType, author string) string {
	var runtimeSection, promptSection string

	switch runtimeType {
	case "prompt":
		runtimeSection = `runtime:
  type: prompt
  entry: prompts/system.md`
		promptSection = `
prompts:
  system: |
    你是一个有帮助的AI助手。
    请根据用户的问题提供准确、有用的回答。`

	case "python":
		runtimeSection = `runtime:
  type: python
  entry: agent.py
  python:
    version: "3.11"
    requirements: requirements.txt`

	case "nodejs":
		runtimeSection = `runtime:
  type: nodejs
  entry: index.js
  nodejs:
    version: "20"
    package: package.json`

	case "docker":
		runtimeSection = `runtime:
  type: docker
  docker:
    dockerfile: Dockerfile`

	case "remote":
		runtimeSection = `runtime:
  type: remote
  remote:
    endpoint: https://your-agent.example.com/api
    protocol: http`
	}

	return fmt.Sprintf(`# AgentSpec - %s
version: "1.0.0"

metadata:
  name: %s
  description: %s
  author: %s
  license: MIT
  tags:
    - %s
  category: %s

%s

model:
  provider: openai
  name: gpt-4o
  parameters:
    temperature: 0.7
    max_tokens: 4096

capabilities:
  streaming: true
  multimodal:
    text: true
    image: false
  memory:
    conversation: true
    long_term: false

interface:
  input:
    type: text
  output:
    type: stream
%s
pricing:
  model: free
`, name, name, description, author, category, category, runtimeSection, promptSection)
}
