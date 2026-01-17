package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	apiURL  string
	Version = "0.1.0"
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "agenthub",
	Short: "AgentHub CLI - 智能体的包管理器",
	Long: `AgentHub CLI 是一个用于管理智能体的命令行工具。

使用 AgentHub，你可以：
  - 搜索和发现智能体
  - 下载和运行智能体  
  - 发布你自己的智能体
  - 管理智能体版本

示例:
  agenthub search "code review"     搜索智能体
  agenthub pull user/agent          下载智能体
  agenthub run user/agent           运行智能体
  agenthub push                     发布智能体
  agenthub login                    登录账户`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认 $HOME/.agenthub/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "https://api.agenthub.dev", "API 服务地址")

	viper.BindPFlag("api_url", rootCmd.PersistentFlags().Lookup("api-url"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configDir := home + "/.agenthub"
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			os.MkdirAll(configDir, 0755)
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("AGENTHUB")

	viper.ReadInConfig()
}
