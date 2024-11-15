package upload

import (
	"PanUpload/cmd/flags"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var cloudreveSession string
var sessionLastTime time.Time

const cloudreveSessionKey = "cloudreve-session"

const ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"

func request(method string, url string, res interface{}, beforeRequest req.RequestMiddleware) error {
	r := req.R()
	r.SetSuccessResult(&res)
	r.SetCookies(&http.Cookie{Name: cloudreveSessionKey, Value: cloudreveSession})
	if beforeRequest != nil {
		err := beforeRequest(r.GetClient(), r)
		if err != nil {
			return err
		}
	}

	resp, err := r.Send(method, url)
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return errors.New(resp.String())
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == cloudreveSessionKey {
			cloudreveSession = cookie.Value
			refreshCloudreveSession()
		}
	}

	return nil
}

// exists checks if file exists.
func exists(path string) bool {
	stat, err := os.Stat(path)
	if err == nil {
		return !stat.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func initCloudreveSession() {
	session, err := os.Open(flags.SessionPath)
	if err != nil {
		fmt.Println("open session file error:", flags.SessionPath, err)
		os.Exit(1)
	}
	defer session.Close()

	// 获取session 刷新时间
	stat, _ := session.Stat()
	sessionLastTime = stat.ModTime()

	data, err := io.ReadAll(session)
	if err != nil {
		fmt.Println("Error reading session file:", flags.SessionPath, err)
		os.Exit(1)
	}

	if len(data) == 0 {
		fmt.Println("Error reading session file:", flags.SessionPath, ",file is empty")
		os.Exit(1)
	}
	cloudreveSession = string(data)
	fmt.Println("read session file:", flags.SessionPath, " success, session:", cloudreveSession)
}

func refreshCloudreveSession() {
	session, err := os.OpenFile(flags.SessionPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("open session file error:", flags.SessionPath, err)
		os.Exit(1)
	}
	defer session.Close()

	// 写入数据到文件
	_, err = session.WriteString(cloudreveSession)
	if err != nil {
		fmt.Println("Error writing to session file:", flags.SessionPath, err)
		os.Exit(1)
	}
	fmt.Println("Success refresh session file:", flags.SessionPath)
}

type ProgressReader struct {
	io.ReadCloser
	totalSize int64
	uploaded  int64
	startTime time.Time
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.ReadCloser.Read(p)
	if n > 0 {
		pr.uploaded += int64(n)
		elapsed := time.Since(pr.startTime).Seconds()
		var speed float64
		if elapsed == 0 {
			speed = float64(pr.uploaded) / 1024
		} else {
			speed = float64(pr.uploaded) / 1024 / elapsed // KB/s
		}

		// 计算进度百分比
		percent := float64(pr.uploaded) / float64(pr.totalSize) * 100
		fmt.Printf("\ruploading: %.2f%% (%d/%d bytes, %.2f KB/s)", percent, pr.uploaded, pr.totalSize, speed)
		// 相等即已经处理完毕
		if pr.uploaded == pr.totalSize {
			fmt.Println()
		}
	}
	return n, err
}

func initClient() {
	initCloudreveSession()
	//req.DevMode() // 将包名视为 Client 直接调用，启用开发模式
	client := req.SetCommonHeader("User-Agent", ua).
		SetCommonHeader("Accept", "application/json, text/plain, */*").
		SetTimeout(30 * time.Minute)
	client.GetTransport().
		WrapRoundTripFunc(func(rt http.RoundTripper) req.HttpRoundTripFunc {
			return func(req *http.Request) (resp *http.Response, err error) {
				// 由于内容长度部分是由后台计算的，所以这里需要手动设置,http默认会过滤掉header.reqWriteExcludeHeader
				if req.ContentLength <= 0 {
					if req.Header.Get("Content-Length") != "" {
						req.ContentLength, _ = strconv.ParseInt(req.Header.Get("Content-Length"), 10, 64)
					}
				}
				return rt.RoundTrip(req)
			}
		})

}
