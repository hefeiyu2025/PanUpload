package cmd

import (
	"PanUpload/cmd/flags"
	"PanUpload/upload"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "pupload",
	Short: "一个huang111.com 的上传工具",
	Run: func(cmd *cobra.Command, args []string) {
		initPath := ""
		if len(args) > 0 {
			initPath = args[0]
		}
		upload.StartUpload(initPath)
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&flags.RemotePath, "remote-path", "r", "/", "云盘目录")
	RootCmd.PersistentFlags().StringVarP(&flags.SessionPath, "session-path", "s", "session.txt", "云盘session")
	RootCmd.PersistentFlags().StringVarP(&flags.CachePath, "cache-path", "c", "./cache", "云盘session")
	RootCmd.PersistentFlags().BoolVarP(&flags.Debug, "debug", "d", false, "开启debug模式")
}
