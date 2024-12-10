package core

import (
	"PanUpload/internal"
	"fmt"
	"github.com/caiguanhao/opencc/configs/t2s"
	client "github.com/hefeiyu2025/pan-client"
	"github.com/hefeiyu2025/pan-client/pan"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func moveInit() (pan.Driver, pan.Driver) {
	internal.InitConfig()
	c, err := client.GetClient(pan.Cloudreve)
	if err != nil {
		panic(err)
	}
	q, err := client.GetClient(pan.Quark)
	if err != nil {
		panic(err)
	}
	return c, q
}

func StartMove() {
	fileChan := make(chan string)
	doneChan := make(chan struct{})
	moveConfig := internal.Config.Move
	cloudreve, quark := moveInit()
	go func() {
		for {
			select {
			case file := <-fileChan:
				if file == "" {
					fmt.Println("接收到全部下载完成")
					close(doneChan)
					return
				}
				absolutePath, err := filepath.Abs(moveConfig.TmpPath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				relativePath, err := filepath.Rel(absolutePath, filepath.Dir(file))
				if err != nil {
					fmt.Println(err)
					continue
				}
				relativePath = strings.ReplaceAll(relativePath, "\\", "/")
				if relativePath == "." {
					relativePath = "/"
				}
				err = quark.UploadFile(pan.UploadFileReq{
					LocalFile:  file,
					RemotePath: relativePath,
					Resumable:  true,
					SuccessDel: true,
					RemotePathTransfer: func(remote string) string {
						return rename(moveConfig, remote)
					},
					RemoteNameTransfer: func(remote string) string {
						return rename(moveConfig, remote)
					},
				})
				if err != nil {
					fmt.Println(err)
					continue
				}
				empty, e := internal.IsEmptyDir(absolutePath)
				if e != nil {
					fmt.Println(e)
					continue
				}
				if empty {
					_ = os.Remove(absolutePath)
				}
			}
		}
	}()
	go func() {
		remotePath := rename(moveConfig, strings.TrimLeft(moveConfig.RemotePath, "/"))
		objs, err := quark.List(pan.ListReq{
			Reload: true,
			Dir: &pan.PanObj{
				Name: remotePath,
				Path: "/",
				Type: "dir",
			},
		})
		ignoreFiles := make([]string, 0)
		ignorePaths := make([]string, 0)
		if err == nil {
			for _, obj := range objs {
				if obj.Type == "dir" {
					ignorePaths = append(ignorePaths, obj.Name)
				} else {
					ignoreFiles = append(ignoreFiles, obj.Name)
				}
			}
		} else {
			fmt.Println(err)
		}

		err = cloudreve.DownloadPath(pan.DownloadPathReq{
			RemotePath: &pan.PanObj{
				Name: strings.TrimLeft(moveConfig.RemotePath, "/"),
				Path: "/",
				Type: "dir",
			},
			SkipFileErr: true,
			LocalPath:   filepath.Join(moveConfig.TmpPath, moveConfig.RemotePath),
			Concurrency: moveConfig.DownloadThread,
			ChunkSize:   moveConfig.DownloadChunkSize,
			OverCover:   false,
			IgnorePaths: ignorePaths,
			IgnoreFiles: ignoreFiles,
			RemoteNameTransfer: func(remote string) string {
				return rename(moveConfig, remote)
			},
			DownloadCallback: func(localPath, localFile string) {
				fileChan <- localFile
			},
		})
		if err != nil {
			fmt.Println(err)
			close(doneChan)
		} else {
			close(fileChan)
		}
	}()

	select {
	case <-doneChan:
		return
	}

}

func rename(moveConfig *internal.MoveConfig, remote string) string {
	newRemote := remote
	for _, removeStr := range moveConfig.RemoveStr {
		newRemote = strings.ReplaceAll(newRemote, removeStr, "")
	}
	// 使用正则表达式替换字符串
	if moveConfig.RemoveReg != "" {
		re := regexp.MustCompile(moveConfig.RemoveReg)
		newRemote = re.ReplaceAllString(newRemote, "")
	}

	newRemote = strings.TrimSpace(newRemote)
	newRemote = t2s.Dicts.Convert(newRemote)

	return newRemote
}
