package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录 AgentHub 账户",
	Long:  `使用用户名/邮箱和密码登录 AgentHub 账户。`,
	Run:   runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("用户名或邮箱: ")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	fmt.Print("密码: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\n读取密码失败")
		os.Exit(1)
	}
	password := string(passwordBytes)
	fmt.Println()

	// 发送登录请求
	apiURL := viper.GetString("api_url")
	reqBody, _ := json.Marshal(map[string]string{
		"login":    login,
		"password": password,
	})

	resp, err := http.Post(apiURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("登录失败: 用户名或密码错误")
		os.Exit(1)
	}

	var result struct {
		User struct {
			Username string `json:"username"`
		} `json:"user"`
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		os.Exit(1)
	}

	// 保存 token
	viper.Set("token", result.Token)
	viper.Set("username", result.User.Username)

	home, _ := os.UserHomeDir()
	configPath := home + "/.agenthub/config.yaml"
	viper.WriteConfigAs(configPath)

	fmt.Printf("✓ 登录成功！欢迎回来，%s\n", result.User.Username)
}

// logoutCmd 登出命令
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "登出当前账户",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("token", "")
		viper.Set("username", "")

		home, _ := os.UserHomeDir()
		configPath := home + "/.agenthub/config.yaml"
		viper.WriteConfigAs(configPath)

		fmt.Println("✓ 已登出")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

// whoamiCmd 查看当前用户
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "显示当前登录用户",
	Run: func(cmd *cobra.Command, args []string) {
		username := viper.GetString("username")
		if username == "" {
			fmt.Println("未登录。使用 'agenthub login' 登录。")
			return
		}
		fmt.Printf("当前登录用户: %s\n", username)
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
