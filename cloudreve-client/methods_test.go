package cloudreve

import (
	"fmt"
	"testing"
)

var cloudreveSession = "MTczMDY4MjkxN3xOd3dBTkVGTVZsVXpWelZVV1RkRFJWUkVRa3BJVVZNM1NsTkJTRVUxU1ZrM1QwSkpWbFJSVjBsVVNqSkVXakpMUzFsTVZ6TkhNMUU9fGn6FJQILddlrz5oIVWWphbWmtexz6f3zCc_zGiTjmq2"

func beforeClient() *CloudreveClient {
	return NewClient("https://pan.huang1111.cn", cloudreveSession)
}

func TestConfig(t *testing.T) {
	client := beforeClient()
	resp, err := client.Config()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileUploadGetUploadSession(t *testing.T) {
	client := beforeClient()

	resp, err := client.FileUploadGetUploadSession(CreateUploadSessionReq{
		// TODO
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileUploadDeleteUploadSession(t *testing.T) {
	client := beforeClient()

	resp, err := client.FileUploadDeleteUploadSession("")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileUploadDeleteAllUploadSession(t *testing.T) {
	client := beforeClient()

	resp, err := client.FileUploadDeleteAllUploadSession()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileCreateFile(t *testing.T) {
	client := beforeClient()

	resp, err := client.FileCreateFile("/11/33.txt")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileCreateDownloadSession(t *testing.T) {
	client := beforeClient()

	resp, err := client.FileCreateDownloadSession("mqoRMnTX")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

//func TestFilePreview(t *testing.T) {
//	client := beforeClient()
//
//	resp, err := client.FilePreview("mqoRMnTX")
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//	}
//	fmt.Println(resp)
//}

func TestFileGetSource(t *testing.T) {
	client := beforeClient()
	dirs := make([]string, 0)
	dirs = append(dirs, "")
	items := make([]string, 0)
	items = append(items, "mqoRMnTX")

	resp, err := client.FileGetSource(ItemReq{
		Item: Item{
			Dirs:  dirs,
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileArchive(t *testing.T) {
	client := beforeClient()
	dirs := make([]string, 0)
	dirs = append(dirs, "DVBmxvCo")
	items := make([]string, 0)
	items = append(items, "")

	resp, err := client.FileArchive(ItemReq{
		Item: Item{
			Dirs:  dirs,
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestCreateDirectory(t *testing.T) {
	client := beforeClient()
	resp, err := client.CreateDirectory("/11")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestListDirectory(t *testing.T) {
	client := beforeClient()
	resp, err := client.ListDirectory("/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestObjectDelete(t *testing.T) {
	client := beforeClient()

	dirs := make([]string, 0)
	dirs = append(dirs, "")
	items := make([]string, 0)
	items = append(items, "6KZbgYC2")

	resp, err := client.ObjectDelete(ItemReq{
		Item: Item{
			Dirs:  dirs,
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestObjectMove(t *testing.T) {
	client := beforeClient()

	dirs := make([]string, 0)
	dirs = append(dirs, "")
	items := make([]string, 0)
	items = append(items, "6KZbgYC2")

	resp, err := client.ObjectMove(ItemMoveReq{
		SrcDir: "/",
		Dst:    "/11",
		Src: Item{
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestObjectCopy(t *testing.T) {
	client := beforeClient()

	dirs := make([]string, 0)
	dirs = append(dirs, "")
	items := make([]string, 0)
	items = append(items, "D1lD4aHB")

	resp, err := client.ObjectCopy(ItemMoveReq{
		SrcDir: "/11",
		Dst:    "/22",
		Src: Item{
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestObjectRename(t *testing.T) {
	client := beforeClient()

	dirs := make([]string, 0)
	dirs = append(dirs, "")
	items := make([]string, 0)
	items = append(items, "6KZLBVT2")

	resp, err := client.ObjectRename(ItemRenameReq{
		NewName: "22.txt",
		Src: Item{
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestObjectGetProperty(t *testing.T) {
	client := beforeClient()

	resp, err := client.ObjectGetProperty(ItemPropertyReq{
		Id:        "E7qQZnHK",
		IsFolder:  true,
		TraceRoot: true,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}
