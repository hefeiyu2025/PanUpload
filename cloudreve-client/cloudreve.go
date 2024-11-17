package cloudreve

// 注意参考 https://github.com/cloudreve/Cloudreve.git

type CloudreveClient struct {
	cloudreveUrl     string
	cloudreveSession string
}

func NewClient(cloudreveUrl, cloudreveSession string) *CloudreveClient {
	return &CloudreveClient{
		cloudreveUrl:     cloudreveUrl,
		cloudreveSession: cloudreveSession}
}

func NewClientWithLogin(cloudreveUrl, username, password string) *CloudreveClient {
	// TODO 登录获取session
	return &CloudreveClient{
		cloudreveUrl: cloudreveUrl,
	}

}
