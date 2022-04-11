package api

type RespCommon struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var (
	Success = RespCommon{
		Code: 0,
		Msg:  "ok",
	}
	DefaultErr = RespCommon{
		Code: 10001,
		Msg:  "error",
	}

	UserFormEmpty = RespCommon{
		Code: 11001,
		Msg:  "form params exist empty",
	}

	UserFormIllegal = RespCommon{
		Code: 11002,
		Msg:  "form params are illegal",
	}
	UserEmailNotExist = RespCommon{
		Code: 11003,
		Msg:  "user email is not exist",
	}
	UserNotActive = RespCommon{
		Code: 11004,
		Msg:  "user email is not active",
	}
	UserPswdErr = RespCommon{
		Code: 11005,
		Msg:  "user password is wrong",
	}
)
