package upload

import (
	"PanUpload/cmd/flags"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const reqPrefix = "https://pan.huang1111.cn/api/v3"

// 刷新session
func config() error {
	var r Resp[any]
	err := request(http.MethodGet, reqPrefix+"/site/config", &r, nil)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return fmt.Errorf("code: %d, msg: %s", r.Code, r.Msg)
	}
	return nil
}

// 列出目录，获取policy.id
func directory() (*DirectoryResp, error) {
	var r Resp[DirectoryResp]
	err := request(http.MethodGet, reqPrefix+"/directory%2F", &r, nil)
	if err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("code: %d, msg: %s", r.Code, r.Msg)
	}
	return r.Data, nil
}

func preUpload(uploadBody *USessionReq) (*USessionInfo, error) {
	var r Resp[USessionInfo]
	err := request(http.MethodPut, reqPrefix+"/file/upload", &r, func(c *resty.Client, r *resty.Request) error {
		r.SetBody(uploadBody)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if r.Code != 0 {
		if r.Code == 40004 {
			return nil, ObjectExistError{
				Message: "文件已存在",
			}
		}
		return nil, fmt.Errorf("code: %d, msg: %s", r.Code, r.Msg)
	}
	return r.Data, nil
}

func uploading(url string, chunk *ChunkData) error {
	// 计算chunk数量
	var result map[string]any
	err := request(http.MethodPut, url, &result, func(c *resty.Client, r *resty.Request) error {
		r.SetContentLength(true)
		r.SetHeader("Content-Type", "application/octet-stream")
		r.SetHeader("Content-Range", "bytes "+strconv.Itoa(chunk.StartSize)+"-"+strconv.Itoa(chunk.EndSize)+"/"+strconv.Itoa(chunk.TotalSize))
		r.SetBody(chunk.Buf)
		return nil
	})
	return err
}

func finishUpload(sessionId string) error {
	var result Resp[string]
	err := request(http.MethodPost, reqPrefix+"/callback/onedrive/finish/"+sessionId, result, func(c *resty.Client, req *resty.Request) error {
		req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		req.SetBody("{}")
		return nil
	})
	if err == nil && result.Code == 0 {
		var key string
		err = getCache(sessionId, &key)
		if err == nil {
			err = delCache(key)
			if err == nil {
				err = delCache(sessionId)
			}
		}
	}
	return err
}

type ChunkConsumer func(chunk *ChunkData) error

func chunkSplit(file *os.File, uploadedChunkNum, chunkSize int, chunkKey string, consumer ChunkConsumer) error {
	var buf []byte
	var chunk int
	stat, err := file.Stat()
	totalSize := int(stat.Size())
	chunkNum := (totalSize / chunkSize) + 1
	fmt.Printf("split chunk total size: %d, num:%d \n", totalSize, chunkNum)
	if uploadedChunkNum >= 0 {
		chunk = uploadedChunkNum + 1
		if chunk > chunkNum {
			fmt.Println("file chunk is uploaded before")
			return nil
		}
		// 将文件指针移动到指定的分片位置
		ret, _ := file.Seek(int64(chunk*chunkSize), 0)
		if ret != 0 {
			fmt.Println("file seek to " + strconv.FormatInt(ret, 10))
		}
	}
	for {
		var n int
		buf = make([]byte, chunkSize)
		n, err = io.ReadAtLeast(file, buf, chunkSize)
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			if err == io.EOF {
				break
			}
			return err
		}

		if n == 0 {
			break
		}
		buf = buf[:n]
		startSize := chunk * chunkSize
		endSize := (chunk * chunkSize) + n - 1

		chunkData := &ChunkData{
			StartSize: startSize,
			EndSize:   endSize,
			ChunkSize: chunkSize,
			TotalSize: totalSize,
			ChunkNum:  chunk,
			Buf:       buf,
		}

		percent := float64(endSize+1) / float64(totalSize) * 100
		fmt.Println("start upload chunk: ", chunk+1)
		err := consumer(chunkData)
		if err != nil {
			return err
		}
		fmt.Printf("success upload chunk: %d , %.2f%% \n", chunk+1, percent)
		// 设置分片缓存
		_ = setCache(chunkKey, strconv.Itoa(chunk))
		chunk++
	}
	return nil
}

func preUploadCache(fileInfo os.FileInfo, policyId, path string) (*USessionInfo, error) {
	reqBody := &USessionReq{
		Path:         path,
		Size:         fileInfo.Size(),
		Name:         fileInfo.Name(),
		PolicyId:     policyId,
		LastModified: fileInfo.ModTime().UnixMilli(),
	}
	key, err := md5Hash(reqBody)
	if err != nil {
		return nil, err
	}
	var resp USessionInfo
	err = getCache(key, &resp)
	if err != nil || (resp.Expires > 0 && resp.Expires < int(time.Now().Unix())) {
		upload, preErr := preUpload(reqBody)
		if preErr == nil {
			if resp.SessionID != "" {
				_ = delCache(resp.SessionID)
				_ = delCache("chunk_" + resp.SessionID)
			}
			preErr = setCache(key, upload)
			if preErr == nil {
				preErr = setCache(upload.SessionID, key)
				resp = *upload
			}
		}
		if err != nil {
			err = preErr
		}
	}
	return &resp, err
}

func md5Hash(params any) (string, error) {
	// 将结构体序列化为JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Md5 error:", err)
		return "", err
	}

	// 计算JSON数据的MD5
	hash := md5.Sum(jsonData)
	md5Str := hex.EncodeToString(hash[:])

	fmt.Println("MD5 Hash:", md5Str)

	return md5Str, nil
}

func exitByError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func StartUpload(file string) {
	initClient()
	initDiskv()
	// 超过6小时才刷新session
	if time.Now().UnixMicro()-sessionLastTime.UnixMicro() >= 6*time.Hour.Microseconds() {
		err := config()
		exitByError(err)
	}
	directoryResp, err := directory()
	exitByError(err)
	// 默认为当前目录
	root := "./"
	//判断参数是否为空，为空则上传当前目录，否则上传指定文件
	if file != "" {
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Println("file read error,", file, err)
			return
		}
		if fileInfo.IsDir() {
			root = file
		} else {
			err = uploadFile(file, directoryResp, "")
			exitByError(err)
			os.Exit(2)
		}
	}

	cachePathAbs, _ := filepath.Abs(flags.CachePath)
	sessionPathAbs, _ := filepath.Abs(flags.SessionPath)
	// 遍历目录
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // 可以选择如何处理错误
			return nil
		}
		if info.IsDir() {
			pathAbs, _ := filepath.Abs(path)
			if pathAbs == cachePathAbs {
				return filepath.SkipDir
			}
		} else {
			pathAbs, _ := filepath.Abs(path)
			if pathAbs == sessionPathAbs {
				return nil
			}
			// 获取相对于root的相对路径
			relPath, _ := filepath.Rel(root, path)
			relPath = strings.Replace(relPath, "\\", "/", -1)
			relPath = strings.Replace(relPath, info.Name(), "", 1)
			return uploadFile(path, directoryResp, relPath)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func uploadFile(path string, directoryResp *DirectoryResp, relPath string) error {
	startTime := time.Now()
	fmt.Println("file start upload,", path, startTime.Format("2006-01-02 15:04:05"))
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("file upload error,", path, err)
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Println("file read stat error,", path, err)
		return nil
	}
	if info.Size() == 0 {
		fmt.Println("file size is zero, give up,", path)
		return nil
	}
	uSessionInfo, err := preUploadCache(info, directoryResp.Policy.Id, flags.RemotePath+relPath)
	if err != nil {
		if errors.As(err, &ObjectExistError{}) {
			fmt.Println("file is exist,", info.Name())
		} else {
			fmt.Println("file upload error,", info.Name(), err)
		}
		return nil
	}
	cacheChunkNum := "-1"
	chunkKey := "chunk_" + uSessionInfo.SessionID
	_ = getCache(chunkKey, &cacheChunkNum)
	uploadedChunkNum, _ := strconv.Atoi(cacheChunkNum)

	err = chunkSplit(file, uploadedChunkNum, uSessionInfo.ChunkSize, chunkKey, func(chunk *ChunkData) error {
		return uploading(uSessionInfo.UploadURLs[0], chunk)
	})
	if err != nil {
		fmt.Println("file upload error,", info.Name(), err)
		return nil
	}
	err = finishUpload(uSessionInfo.SessionID)
	if err != nil {
		fmt.Println("file finish upload error,", info.Name(), err)
		return nil
	}
	_ = delCache(chunkKey)
	fmt.Println("file success upload,", path, time.Now().Format("2006-01-02 15:04:05"), time.Since(startTime))
	return nil
}
