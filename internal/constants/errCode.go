package constants

type ErrCode int64

const (
	Success    ErrCode = 0
	DefaultErr ErrCode = 10001

	//User
	UserActiveErr ErrCode = 11001
	UserPswdErr   ErrCode = 11002
	UserEmailErr  ErrCode = 11003
)
