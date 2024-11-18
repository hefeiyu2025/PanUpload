package cloudreve

import (
	"fmt"
	"net/url"
)

func (c *CloudreveClient) Config() (*RespData[SiteConfig], error) {
	r := c.sessionClient.R()
	var successResult RespData[SiteConfig]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	//site/config
	response, err := r.Get("/site/config")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	for _, cookie := range response.Cookies() {
		if cookie.Name == "cloudreve-session" {
			c.refreshSession(cookie.Value)
		}
	}
	return &successResult, nil
}

func (c *CloudreveClient) UserStorage() (*RespData[Storage], error) {
	r := c.sessionClient.R()
	var successResult RespData[Storage]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	// /file/upload
	response, err := r.Get("/user/storage")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

//	func (c *CloudreveClient) FileUpload(sessionId string, index int) Resp {
//		return Resp{}
//	}

func (c *CloudreveClient) FileUploadGetUploadSession(req CreateUploadSessionReq) (*RespData[UploadCredential], error) {
	r := c.sessionClient.R()
	var successResult RespData[UploadCredential]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	r.SetBody(req)
	// /file/upload
	response, err := r.Put("/file/upload")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) FileUploadDeleteUploadSession(sessionId string) (*Resp, error) {
	r := c.sessionClient.R()
	var result Resp
	r.SetSuccessResult(&result)
	r.SetErrorResult(&result)
	// /file/upload/{sessionId}
	response, err := r.Delete("/file/upload/" + sessionId)
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	return &result, nil
}

func (c *CloudreveClient) FileUploadDeleteAllUploadSession() (*Resp, error) {
	r := c.sessionClient.R()
	var result Resp
	r.SetSuccessResult(&result)
	r.SetErrorResult(&result)
	// /file/upload
	response, err := r.Delete("/file/upload")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	return &result, nil
}

//func (c *CloudreveClient) FilePutContent(path string) {
//
//}

func (c *CloudreveClient) FileCreateFile(path string) (*Resp, error) {
	r := c.sessionClient.R()
	var result Resp
	r.SetSuccessResult(&result)
	r.SetErrorResult(&result)
	r.SetBody(Json{
		"path": path,
	})
	// /file/create
	response, err := r.Post("/file/create")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	return &result, nil
}

func (c *CloudreveClient) FileCreateDownloadSession(id string) (*RespData[string], error) {
	r := c.sessionClient.R()
	var successResult RespData[string]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	// /file/download
	response, err := r.Put("/file/download/" + id)
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

//func (c *CloudreveClient) FilePreview(id string) (string, error) {
//	r := c.sessionClient.R()
//
//
//	// /file/preview
//	response, err := r.Get("/file/preview/" + id)
//	if err != nil {
//		return "", err
//	}
//
//	return response.String(), nil
//}

//func (c *CloudreveClient) FilePreviewText(id string) {
//
//}

func (c *CloudreveClient) FileGetSource(req ItemReq) (*RespData[[]Sources], error) {
	r := c.sessionClient.R()
	var successResult RespData[[]Sources]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	r.SetBody(req)
	// /file/source
	response, err := r.Post("/file/source")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) FileArchive(req ItemReq) (*RespData[string], error) {
	r := c.sessionClient.R()
	var successResult RespData[string]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	r.SetBody(req)
	// /file/archive
	response, err := r.Post("/file/archive")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}
func (c *CloudreveClient) CreateDirectory(path string) (*Resp, error) {
	r := c.sessionClient.R()
	var result Resp
	r.SetSuccessResult(&result)
	r.SetErrorResult(&result)
	r.SetBody(Json{
		"path": path,
	})
	//directory
	response, err := r.Put("/directory")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	return &result, nil
}

func (c *CloudreveClient) ListDirectory(path string) (*RespData[ObjectList], error) {
	r := c.sessionClient.R()
	var successResult RespData[ObjectList]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	// /directory*path
	response, err := r.Get("/directory" + url.PathEscape(path))
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) ObjectDelete(req ItemReq) (*Resp, error) {
	r := c.sessionClient.R()
	var result Resp
	r.SetSuccessResult(&result)
	r.SetErrorResult(&result)
	r.SetBody(req)
	// object
	response, err := r.Delete("/object")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	return &result, nil
}

func (c *CloudreveClient) ObjectMove(req ItemMoveReq) (*Resp, error) {
	r := c.sessionClient.R()
	var result Resp
	r.SetSuccessResult(&result)
	r.SetErrorResult(&result)
	r.SetBody(req)
	// object
	response, err := r.Patch("/object")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	return &result, nil
}

func (c *CloudreveClient) ObjectCopy(req ItemMoveReq) (*Resp, error) {
	r := c.sessionClient.R()
	var result Resp
	r.SetSuccessResult(&result)
	r.SetErrorResult(&result)
	r.SetBody(req)
	// /object/copy
	response, err := r.Post("/object/copy")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	return &result, nil
}

func (c *CloudreveClient) ObjectRename(req ItemRenameReq) (*Resp, error) {
	r := c.sessionClient.R()
	var result Resp
	r.SetSuccessResult(&result)
	r.SetErrorResult(&result)
	r.SetBody(req)
	// /object/rename
	response, err := r.Post("/object/rename")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	return &result, nil
}

func (c *CloudreveClient) ObjectGetProperty(req ItemPropertyReq) (*RespData[ObjectProps], error) {
	r := c.sessionClient.R()
	var errorResult Resp
	var successResult RespData[ObjectProps]
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	r.SetQueryParamsAnyType(Json{
		"is_folder":  req.IsFolder,
		"trace_root": req.TraceRoot,
	})
	// /object/property/{id}
	response, err := r.Get("/object/property/" + req.Id)
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) ShareCreateShare(req ShareCreateReq) (*RespData[string], error) {
	r := c.sessionClient.R()
	var successResult RespData[string]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	r.SetBody(req)
	// /share
	response, err := r.Post("/share")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) ShareListShare() (*RespData[ShareList], error) {
	r := c.sessionClient.R()
	var successResult RespData[ShareList]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	// /share
	response, err := r.Get("/share")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) ShareUpdateShare(req ShareUpdateReq) (*RespData[string], error) {
	r := c.sessionClient.R()
	var successResult RespData[string]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	r.SetBody(Json{
		"prop":  req.Prop,
		"value": req.Value,
	})
	// /share
	response, err := r.Patch("/share/" + req.Id)
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) ShareDeleteShare(id string) (*Resp, error) {
	r := c.sessionClient.R()
	var result Resp
	r.SetSuccessResult(&result)
	r.SetErrorResult(&result)
	// /share
	response, err := r.Delete("/share/" + id)
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	return &result, nil
}

func (c *CloudreveClient) ShareGetShare(id, password string) (*RespData[Share], error) {
	r := c.defaultClient.R()
	var successResult RespData[Share]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	r.SetBody(Json{
		"password": password,
	})
	// /share
	response, err := r.Get("/share/info/" + id)
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) ShareGetShareDownload(id, path string) (*RespData[string], error) {
	r := c.defaultClient.R()
	var successResult RespData[string]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	r.SetQueryParamsAnyType(Json{
		"path": path,
	})
	// /share
	response, err := r.Put("/share/download/" + id)
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) ShareListSharedFolder(id, path string) (*RespData[ObjectList], error) {
	r := c.defaultClient.R()
	var successResult RespData[ObjectList]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)

	// /share/list/:id/*path
	response, err := r.Put("/share/list/" + id + url.PathEscape(path))
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) ShareSearchSharedFolder(id, keyword, path string, searchType SearchType) (*RespData[ObjectList], error) {
	r := c.defaultClient.R()
	var successResult RespData[ObjectList]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	r.SetQueryParamsAnyType(Json{
		"path": path,
	})
	// /share/search/:id/:type/:keywords
	response, err := r.Get("/share/search/" + id + "/" + string(searchType) + "/" + keyword)
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) ShareSearchShare(req ShareListReq) (*RespData[ShareList], error) {
	r := c.defaultClient.R()
	var successResult RespData[ShareList]
	var errorResult Resp
	r.SetSuccessResult(&successResult)
	r.SetErrorResult(&errorResult)
	r.SetQueryParamsAnyType(Json{
		"page":     req.Page,
		"order_by": req.OrderBy,
		"order":    req.Order,
		"keywords": req.Keywords,
	})
	// /share/search
	response, err := r.Get("/share/search")
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", errorResult.Code, errorResult.Msg)
	}
	return &successResult, nil
}

func (c *CloudreveClient) OneDriveCallback(sessionId string) (*Resp, error) {
	r := c.sessionClient.R()
	var result Resp
	r.SetSuccessResult(&result)
	r.SetErrorResult(&result)
	// /callback/onedrive/finish/:sessionID
	response, err := r.Get("/callback/onedrive/finish/" + sessionId)
	if err != nil {
		return nil, err
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	return &result, nil
}

//func (c *CloudreveClient) OneDriveUpload(req OneDriveUploadReq) (*Resp, error) {
//
//}
