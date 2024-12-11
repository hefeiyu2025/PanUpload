package main

import (
	"PanUpload/core"
	"fmt"
	client "github.com/hefeiyu2025/pan-client"
	"os"
	"os/signal"
	"syscall"
)

// getPidFile 获取或创建PID文件
func getPidFile() string {
	// 这里使用当前程序的名称作为PID文件的名称
	return fmt.Sprintf("%s.pid", os.Args[0])
}

// writePid 将当前进程的PID写入PID文件
func writePid(pidFile string) error {
	// 获取当前进程的PID
	pid := os.Getpid()
	return os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644)
}

// removePid 移除PID文件
func removePid(pidFile string) {
	err := os.Remove(pidFile)
	if err != nil {
		fmt.Printf("remove pid file error: %s\n", err)
	}
}

// checkPid 检查PID文件是否存在，如果存在则程序已经运行
func checkPid(pidFile string) bool {
	_, err := os.Stat(pidFile)
	return !os.IsNotExist(err)
}

func main() {

	pidFile := getPidFile()

	// 检查PID文件是否存在
	if checkPid(pidFile) {
		fmt.Println("程序已在运行中...")
		os.Exit(1)
	}

	// 写入PID到文件
	if err := writePid(pidFile); err != nil {
		fmt.Printf("写入PID失败: %v\n", err)
		os.Exit(1)
	}

	// 设置信号处理器，以便在程序退出时移除PID文件
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		removePid(pidFile)
		client.GracefulExist()
		os.Exit(1)
	}()

	//cmd.Execute()
	defer func() {
		removePid(pidFile)
		client.GracefulExist()
	}()
	core.StartMove()
}
