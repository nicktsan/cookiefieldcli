package loginInterface

import (
	"cookiefieldcli/cmd/login/loginResponse"
)

type ILogin interface {
	GetDeviceCode() (loginResponse.LResponse, error)
	GetRequestToken(loginResponse.LResponse)
}
