package api

import (
	"encoding/json"
	"net/url"

	"github.com/heroicyang/wechat-crypter"
	"github.com/heroicyang/wechat-qy/base"
)

// 企业号相关接口的 API 接口地址
const (
	FetchTokenURI       = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	CreateMenuURI       = "https://qyapi.weixin.qq.com/cgi-bin/menu/create"
	DeleteMenuURI       = "https://qyapi.weixin.qq.com/cgi-bin/menu/delete"
	GetMenuURI          = "https://qyapi.weixin.qq.com/cgi-bin/menu/get"
	CreateDepartmentURI = "https://qyapi.weixin.qq.com/cgi-bin/department/create"
	UpdateDepartmentURI = "https://qyapi.weixin.qq.com/cgi-bin/department/update"
	DeleteDepartmentURI = "https://qyapi.weixin.qq.com/cgi-bin/department/delete"
	ListDepartmentURI   = "https://qyapi.weixin.qq.com/cgi-bin/department/list"
	CreateUserURI       = "https://qyapi.weixin.qq.com/cgi-bin/user/create"
	UpdateUserURI       = "https://qyapi.weixin.qq.com/cgi-bin/user/update"
	DeleteUserURI       = "https://qyapi.weixin.qq.com/cgi-bin/user/delete"
	BatchDeleteUserURI  = "https://qyapi.weixin.qq.com/cgi-bin/user/batchdelete"
	GetUserURI          = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	ListSimpleUserURI   = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist"
	ListUserURI         = "https://qyapi.weixin.qq.com/cgi-bin/user/list"
	InviteUserURI       = "https://qyapi.weixin.qq.com/cgi-bin/invite/send"
)

// API 封装了企业号相关的接口操作
type API struct {
	corpSecret string
	CorpID     string
	MsgCrypter crypter.MessageCrypter
	Client     *base.Client
	Tokener    base.Tokener
}

// New 方法创建 API 实例
func New(corpID, corpSecret, token, encodingAESKey string) *API {
	msgCrypter, _ := crypter.NewMessageCrypter(token, encodingAESKey, corpID)

	api := &API{
		corpSecret: corpSecret,
		CorpID:     corpID,
		MsgCrypter: msgCrypter,
	}

	api.Client = base.NewClient(api)
	api.Tokener = NewTokener(api)

	return api
}

// Retriable 方法实现了 API 在发起请求遇到 token 错误时，先刷新 token 然后再次发起请求的逻辑
func (a *API) Retriable(body []byte) (bool, error) {
	result := &base.Error{}
	if err := json.Unmarshal(body, result); err != nil {
		return false, err
	}

	switch result.ErrCode {
	case base.ErrCodeOk:
		return false, nil
	case base.ErrCodeTokenInvalid, base.ErrCodeTokenTimeout:
		if _, err := a.Tokener.RefreshToken(); err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, result
	}
}

// FetchToken 方法使用企业管理组的密钥向 API 服务器获取企业号的令牌信息
func (a *API) FetchToken() (token string, expiresIn int64, err error) {
	qs := make(url.Values)
	qs.Add("corpid", a.CorpID)
	qs.Add("corpsecret", a.corpSecret)

	url := FetchTokenURI + "?" + qs.Encode()

	body, err := a.Client.GetJSON(url)
	if err != nil {
		return
	}

	result := &TokenInfo{}
	if err = json.Unmarshal(body, result); err != nil {
		return
	}

	token = result.Token
	expiresIn = result.ExpiresIn

	return
}