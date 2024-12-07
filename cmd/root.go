package cmd

import (
	"PanUpload/upload"
	"fmt"
	"github.com/spf13/cobra"
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
		panic(err)
	}
}
