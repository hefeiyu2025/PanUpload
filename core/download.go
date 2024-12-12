package core

import (
	"PanUpload/internal"
	"github.com/caiguanhao/opencc/configs/t2s"
	client "github.com/hefeiyu2025/pan-client"
	"github.com/hefeiyu2025/pan-client/pan"
	"path/filepath"
	"regexp"
	"strings"
)

func downloadInit() pan.Driver {
	internal.InitConfig()
	c, err := client.GetClient(pan.DriverType(internal.Config.Download.DownloadClient))
	if err != nil {
		panic(err)
	}
	return c
}

func StartDownload() {
	driver := downloadInit()
	downloadConfig := internal.Config.Download
	remoteTransfer := func(remote string) string {
		newRemote := remote
		for _, removeStr := range downloadConfig.RemoveStr {
			newRemote = strings.ReplaceAll(newRemote, removeStr, "")
		}
		// 使用正则表达式替换字符串
		if downloadConfig.RemoveReg != "" {
			re := regexp.MustCompile(downloadConfig.RemoveReg)
			newRemote = re.ReplaceAllString(newRemote, "")
		}

		newRemote = strings.TrimSpace(newRemote)
		newRemote = t2s.Dicts.Convert(newRemote)

		return newRemote
	}
	paths := downloadConfig.RemotePath
	for _, path := range paths {
		err := driver.DownloadPath(pan.DownloadPathReq{
			RemotePath: &pan.PanObj{
				Name: strings.TrimLeft(path, "/"),
				Path: "/",
				Type: "dir",
			},
			SkipFileErr:        true,
			LocalPath:          filepath.Join(downloadConfig.LocalPath, path),
			Concurrency:        downloadConfig.DownloadThread,
			ChunkSize:          downloadConfig.DownloadChunkSize,
			RemoteNameTransfer: remoteTransfer,
		})

		if err != nil {
			panic(err)
		}
	}

}
