package cloudreve

import "time"

// Response 基础序列化器
type Response struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Error string `json:"error,omitempty"`
}

type ResponseData[T any] struct {
	Response
	Data T `json:"data,omitempty"`
}

type SiteConfig struct {
	SiteName             string   `json:"title"`
	LoginCaptcha         bool     `json:"loginCaptcha"`
	RegCaptcha           bool     `json:"regCaptcha"`
	ForgetCaptcha        bool     `json:"forgetCaptcha"`
	EmailActive          bool     `json:"emailActive"`
	Themes               string   `json:"themes"`
	DefaultTheme         string   `json:"defaultTheme"`
	HomepageViewMethod   string   `json:"home_view_method"`
	ShareViewMethod      string   `json:"share_view_method"`
	Authn                bool     `json:"authn"`
	User                 User     `json:"user"`
	ReCaptchaKey         string   `json:"captcha_ReCaptchaKey"`
	CaptchaType          string   `json:"captcha_type"`
	TCaptchaCaptchaAppId string   `json:"tcaptcha_captcha_app_id"`
	RegisterEnabled      bool     `json:"registerEnabled"`
	AppPromotion         bool     `json:"app_promotion"`
	WopiExts             []string `json:"wopi_exts"`
}

// User 用户序列化器
type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"user_name"`
	Nickname       string    `json:"nickname"`
	Status         int       `json:"status"`
	Avatar         string    `json:"avatar"`
	CreatedAt      time.Time `json:"created_at"`
	PreferredTheme string    `json:"preferred_theme"`
	Anonymous      bool      `json:"anonymous"`
	Group          group     `json:"group"`
	Tags           []tag     `json:"tags"`
}

type group struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	AllowShare           bool   `json:"allowShare"`
	AllowRemoteDownload  bool   `json:"allowRemoteDownload"`
	AllowArchiveDownload bool   `json:"allowArchiveDownload"`
	ShareDownload        bool   `json:"shareDownload"`
	CompressEnabled      bool   `json:"compress"`
	WebDAVEnabled        bool   `json:"webdav"`
	SourceBatchSize      int    `json:"sourceBatch"`
	AdvanceDelete        bool   `json:"advanceDelete"`
	AllowWebDAVProxy     bool   `json:"allowWebDAVProxy"`
}

type tag struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Color      string `json:"color"`
	Type       int    `json:"type"`
	Expression string `json:"expression"`
}

type storage struct {
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Total uint64 `json:"total"`
}

// CreateUploadSessionReq 获取上传凭证服务
type CreateUploadSessionReq struct {
	Path         string `json:"path"`
	Size         uint64 `json:"size"`
	Name         string `json:"name" `
	PolicyID     string `json:"policy_id" `
	LastModified int64  `json:"last_modified"`
	MimeType     string `json:"mime_type"`
}

// UploadCredential 返回给客户端的上传凭证
type UploadCredential struct {
	SessionID   string   `json:"sessionID"`
	ChunkSize   uint64   `json:"chunkSize"` // 分块大小，0 为部分快
	Expires     int64    `json:"expires"`   // 上传凭证过期时间， Unix 时间戳
	UploadURLs  []string `json:"uploadURLs,omitempty"`
	Credential  string   `json:"credential,omitempty"`
	UploadID    string   `json:"uploadID,omitempty"`
	Callback    string   `json:"callback,omitempty"` // 回调地址
	Path        string   `json:"path,omitempty"`     // 存储路径
	AccessKey   string   `json:"ak,omitempty"`
	KeyTime     string   `json:"keyTime,omitempty"` // COS用有效期
	Policy      string   `json:"policy,omitempty"`
	CompleteURL string   `json:"completeURL,omitempty"`
}
