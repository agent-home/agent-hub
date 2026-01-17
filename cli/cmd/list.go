package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出智能体",
	Long: `列出热门智能体或特定分类的智能体。

示例:
  agenthub list                    # 列出热门智能体
  agenthub list --category coding  # 列出编程类智能体
  agenthub list --mine             # 列出我的智能体`,
	Run: runList,
}

var (
	listCategory string
	listMine     bool
	listPage     int
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listCategory, "category", "c", "", "按分类筛选")
	listCmd.Flags().BoolVarP(&listMine, "mine", "m", false, "只显示我的智能体")
	listCmd.Flags().IntVarP(&listPage, "page", "p", 1, "页码")
}

func runList(cmd *cobra.Command, args []string) {
	apiURL := viper.GetString("api_url")

	var reqURL string

	if listMine {
		username := viper.GetString("username")
		if username == "" {
			fmt.Println("请先登录: agenthub login")
			os.Exit(1)
		}
		reqURL = fmt.Sprintf("%s/api/v1/users/%s/agents", apiURL, username)
	} else {
		reqURL = fmt.Sprintf("%s/api/v1/agents?page=%d&page_size=20", apiURL, listPage)
		if listCategory != "" {
			reqURL += "&category=" + listCategory
		}
	}

	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result struct {
		Agents []AgentInfo `json:"agents"`
		Total  int64       `json:"total"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		os.Exit(1)
	}

	if len(result.Agents) == 0 {
		if listMine {
			fmt.Println("你还没有发布任何智能体")
			fmt.Println("使用 'agenthub init' 创建，'agenthub push' 发布")
		} else {
			fmt.Println("暂无智能体")
		}
		return
	}

	if listMine {
		fmt.Println("我的智能体:")
	} else if listCategory != "" {
		fmt.Printf("%s 分类的智能体:\n", listCategory)
	} else {
		fmt.Println("热门智能体:")
	}
	fmt.Println()

	for _, agent := range result.Agents {
		fullName := agent.Namespace + "/" + agent.Name
		if agent.FullName != "" {
			fullName = agent.FullName
		}

		fmt.Printf("📦 %s\n", fullName)
		if agent.Description != "" {
			desc := agent.Description
			if len(desc) > 70 {
				desc = desc[:67] + "..."
			}
			fmt.Printf("   %s\n", desc)
		}
		fmt.Printf("   ⬇️ %d  ❤️ %d\n", agent.Downloads, agent.Likes)
		fmt.Println()
	}

	if !listMine && result.Total > 20 {
		fmt.Printf("共 %d 个智能体，使用 --page 查看更多\n", result.Total)
	}
}

// categoriesCmd 列出所有分类
var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "列出所有分类",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := viper.GetString("api_url")

		resp, err := http.Get(apiURL + "/api/v1/categories")
		if err != nil {
			fmt.Printf("请求失败: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		var result struct {
			Categories []struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Count int    `json:"count"`
			} `json:"categories"`
		}

		json.NewDecoder(resp.Body).Decode(&result)

		fmt.Println("可用分类:")
		fmt.Println()
		for _, cat := range result.Categories {
			fmt.Printf("  %-15s %s\n", cat.ID, cat.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
}
