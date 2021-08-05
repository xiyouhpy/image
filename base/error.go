// Package base 提供自定义错误的方法
package base

// 错误号code信息
const (
	// ErrCodeSuccess 通用错误号，错误号范围 [0, 10000)
	ErrCodeSuccess            = 0  // 请求成功
	ErrCodeUnknownError       = 1  // 未知错误
	ErrCodeParamsError        = 2  // 参数错误
	ErrCodeConfError          = 3  // 获取配置失败
	ErrCodeJSONMarshalError   = 4  // json编码失败
	ErrCodeJSONUnMarshalError = 5  // json解码失败
	ErrCodeImgEncodeError     = 6  // img编码错误
	ErrCodeImgDecodeError     = 7  // img解码错误
	ErrCodeRedisError         = 8  // 调用 redis 失败
	ErrCodeMysqlError         = 9  // 调用 mysql 失败
	ErrCodeServiceError       = 10 // 服务错误
	ErrCodeSqlError           = 11 // sql错误
	ErrCodeDownloadError      = 12 // 下载失败
	ErrCodeUploadError        = 13 // 上传失败

	// 业务错误号，错误号范围 [10000, 999999)；前两位表示属于哪个业务，后四位表示该业务的细分错误号
	// todo...
)

// 错误号msg信息
const (
	ErrMsgSuccess            = "success"
	ErrMsgUnknownError       = "Unknown err"
	ErrMsgParamsError        = "params err"
	ErrMsgConfError          = "conf err"
	ErrMsgJSONMarshalError   = "json marshal err"
	ErrMsgJSONUnMarshalError = "json unmarshal err"
	ErrMsgImgEncodeError     = "img encode err"
	ErrMsgImgDecodeError     = "img decode err"
	ErrMsgRedisError         = "redis err"
	ErrMsgMysqlError         = "mysql err"
	ErrMsgServiceError       = "service err"
	ErrMsgSqlError           = "sql err"
	ErrMsgDownloadError      = "download err"
	ErrMsgUploadError        = "upload err"
)

// 错误号和错误信息映射关系
var (
	ErrSuccess            = NewErrno(ErrCodeSuccess, ErrMsgSuccess)
	ErrUnknownError       = NewErrno(ErrCodeUnknownError, ErrMsgUnknownError)
	ErrParamsError        = NewErrno(ErrCodeParamsError, ErrMsgParamsError)
	ErrConfError          = NewErrno(ErrCodeConfError, ErrMsgConfError)
	ErrJSONMarshalError   = NewErrno(ErrCodeJSONMarshalError, ErrMsgJSONMarshalError)
	ErrJSONUnMarshalError = NewErrno(ErrCodeJSONUnMarshalError, ErrMsgJSONUnMarshalError)
	ErrImgEncodeError     = NewErrno(ErrCodeImgEncodeError, ErrMsgImgEncodeError)
	ErrImgDecodeError     = NewErrno(ErrCodeImgDecodeError, ErrMsgImgDecodeError)
	ErrRedisError         = NewErrno(ErrCodeRedisError, ErrMsgRedisError)
	ErrMysqlError         = NewErrno(ErrCodeMysqlError, ErrMsgMysqlError)
	ErrServiceError       = NewErrno(ErrCodeServiceError, ErrMsgServiceError)
	ErrSqlError           = NewErrno(ErrCodeSqlError, ErrMsgSqlError)
	ErrDownloadError      = NewErrno(ErrCodeDownloadError, ErrMsgDownloadError)
	ErrUploadError        = NewErrno(ErrCodeUploadError, ErrMsgUploadError)
)

// ErrCode 错误号类型定义
type ErrCode uint32

// Errno 自定义错误号结构
type myErrno struct {
	code ErrCode
	msg  string
}

// Error 自定义错误号结构
type Error interface {
	error
	Errno() ErrCode
}

// NewErrno 错误号对象
func NewErrno(code ErrCode, msg string) Error {
	return &myErrno{code, msg}
}

// ErrMsg 获取错误描述方法
func (err *myErrno) Error() string {
	return err.msg
}

// Errno 获取错误号方法
func (err *myErrno) Errno() ErrCode {
	return err.code
}
