package cloudreve

// site/config
func Config() ResponseData[SiteConfig] {
	return ResponseData[SiteConfig]{}
}

func FileUpload(sessionId string, index int) Response {
	return Response{}
}

func FileUploadGetUploadSession(req CreateUploadSessionReq) ResponseData[UploadCredential] {
	return ResponseData[UploadCredential]{}
}

func FileUploadDeleteUploadSession(sessionId string) Response {

}

func FileUploadDeleteAllUploadSession() Response {

}

func FileUpdate() {

}
