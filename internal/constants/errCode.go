package constants

type ErrCode int64

const (
	Success    ErrCode = 0
	DefaultErr ErrCode = 10001

	//User
	UserEmailErr   ErrCode = 11001
	UserActiveErr  ErrCode = 11002
	UserPswdErr    ErrCode = 11003
	UserEmailExist ErrCode = 11004
	UserSidExist   ErrCode = 11005

	//Doc
	DocTitleExist   ErrCode = 12001
	DocUploadQnyErr ErrCode = 12002
)
