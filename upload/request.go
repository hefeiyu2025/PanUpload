package upload

import (
	"PanUpload/cmd/flags"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"os"
	"time"
)

var client *resty.Client
var cloudreveSession string
var sessionLastTime time.Time

const cloudreveSessionKey = "cloudreve-session"

const ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"

func request(method string, url string, res interface{}, beforeRequest resty.RequestMiddleware) error {
	req := client.R()

	req.SetResult(&res)
	req.SetHeaders(map[string]string{
		"Accept":     "application/json, text/plain, */*",
		"User-Agent": ua,
	})
	req.SetCookie(&http.Cookie{Name: cloudreveSessionKey, Value: cloudreveSession})
	if beforeRequest != nil {
		err := beforeRequest(client, req)
		if err != nil {
			return err
		}
	}
	resp, err := req.Execute(method, url)
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
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

func preRequestHook(c *resty.Client, h *http.Request) error {
	if h.Body != nil {
		size := h.ContentLength
		pr := &ProgressReader{
			ReadCloser: h.Body,
			startTime:  time.Now(),
			totalSize:  size,
		}
		h.Body = pr
	}
	return nil
}

func initClient() {
	initCloudreveSession()
	client = resty.New().
		SetRetryCount(3).
		SetRetryResetReaders(true).
		SetTimeout(30 * time.Minute).
		SetRateLimiter(rate.NewLimiter(rate.Every(time.Second), 30)).
		SetDebug(flags.Debug).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetPreRequestHook(preRequestHook)
}
