package cmd

import (
	"PanUpload/cmd/flags"
	"PanUpload/upload"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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

func init() {
	process, _ := os.Executable()
	processDir := filepath.Dir(process)
	RootCmd.PersistentFlags().StringVarP(&flags.RemotePath, "remote-path", "r", "/", "云盘目录")
	RootCmd.PersistentFlags().StringVarP(&flags.SessionPath, "session-path", "s", processDir+"/session.txt", "云盘session")
	RootCmd.PersistentFlags().StringVarP(&flags.CachePath, "cache-path", "c", processDir+"/cache", "云盘session")
	RootCmd.PersistentFlags().BoolVarP(&flags.Debug, "debug", "d", false, "开启debug模式")
	RootCmd.PersistentFlags().StringVarP(&flags.RemoveStr, "remove-str", "t", "TVBOXNOW,ViuTV", "替换的名称")
	RootCmd.PersistentFlags().StringVarP(&flags.RemoveReg, "remove-reg", "g", "\\(\\d+\\)", "替换的正则表达表达式")
	RootCmd.PersistentFlags().BoolVarP(&flags.Delete, "delete", "o", false, "上传成功是否删除文件")
	RootCmd.PersistentFlags().BoolVarP(&flags.DeleteAllSession, "del-session", "d", false, "上传前是否清空所有session")
	RootCmd.PersistentFlags().StringVarP(&flags.UploadExtensions, "upload-extensions", "e", ".mp4,.avi,.mkv,.flv,.mov,.rmvb,.rm,.ts", "处理的文件后缀")
	RootCmd.PersistentFlags().StringVarP(&flags.IgnorePath, "ignore-path", "i", "云盘缓存文件", "不进行处理的目录")
}
