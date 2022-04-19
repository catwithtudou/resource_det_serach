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
	FormEmptyErr = RespCommon{
		Code: 10002,
		Msg:  "form params exist empty",
	}

	FormIllegalErr = RespCommon{
		Code: 10003,
		Msg:  "form params are illegal",
	}

	UserAuthErr = RespCommon{
		Code: 10004,
		Msg:  "user auth err",
	}

	FormFileErr = RespCommon{
		Code: 10005,
		Msg:  "form file err",
	}

	FileTypeErr = RespCommon{
		Code: 10006,
		Msg:  "file type not supported",
	}

	FileSizeErr = RespCommon{
		Code: 10007,
		Msg:  "file size not supported",
	}

	UserEmailNotExist = RespCommon{
		Code: 11001,
		Msg:  "user email is not exist",
	}
	UserNotActive = RespCommon{
		Code: 11002,
		Msg:  "user email is not active",
	}
	UserPswdErr = RespCommon{
		Code: 11003,
		Msg:  "user password is wrong",
	}
	UserEmailExist = RespCommon{
		Code: 11004,
		Msg:  "user email is exist",
	}
	UserSidExist = RespCommon{
		Code: 11005,
		Msg:  "user sid is exist",
	}

	DocTitleExist = RespCommon{
		Code: 12001,
		Msg:  "user doc title is exist",
	}

	DocUploadQnyErr = RespCommon{
		Code: 12002,
		Msg:  "user upload doc qny err",
	}
)
