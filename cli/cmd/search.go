package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	searchCategory string
	searchLimit    int
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "搜索智能体",
	Long: `搜索 AgentHub 上的智能体。

示例:
  agenthub search "code review"
  agenthub search assistant --category coding
  agenthub search "data analysis" --limit 20`,
	Args: cobra.MinimumNArgs(1),
	Run:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&searchCategory, "category", "c", "", "按分类筛选")
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 10, "结果数量限制")
}

type SearchResult struct {
	Agents []AgentInfo `json:"agents"`
	Total  int64       `json:"total"`
}

type AgentInfo struct {
	FullName    string   `json:"full_name"`
	Name        string   `json:"name"`
	Namespace   string   `json:"namespace"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Downloads   int64    `json:"downloads"`
	Likes       int64    `json:"likes"`
}

func runSearch(cmd *cobra.Command, args []string) {
	query := strings.Join(args, " ")
	apiURL := viper.GetString("api_url")

	// 构建请求 URL
	reqURL := fmt.Sprintf("%s/api/v1/search?q=%s&page_size=%d",
		apiURL, url.QueryEscape(query), searchLimit)

	if searchCategory != "" {
		reqURL += "&category=" + url.QueryEscape(searchCategory)
	}

	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Printf("搜索失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		os.Exit(1)
	}

	if len(result.Agents) == 0 {
		fmt.Printf("未找到匹配 \"%s\" 的智能体\n", query)
		return
	}

	fmt.Printf("找到 %d 个智能体:\n\n", result.Total)

	for _, agent := range result.Agents {
		fullName := agent.Namespace + "/" + agent.Name
		if agent.FullName != "" {
			fullName = agent.FullName
		}

		fmt.Printf("📦 %s\n", fullName)
		if agent.Description != "" {
			desc := agent.Description
			if len(desc) > 80 {
				desc = desc[:77] + "..."
			}
			fmt.Printf("   %s\n", desc)
		}
		fmt.Printf("   ⬇️ %d  ❤️ %d", agent.Downloads, agent.Likes)
		if agent.Category != "" {
			fmt.Printf("  📁 %s", agent.Category)
		}
		fmt.Println()
		fmt.Println()
	}

	fmt.Printf("使用 'agenthub pull <name>' 下载智能体\n")
}
