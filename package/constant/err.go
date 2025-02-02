package constant

// 本错误码参照了HTTP状态码的语义，方便识别错误类型

//分类	分类描述
//1**	信息，服务器收到请求，需要请求者继续执行操作
//2**	成功，操作被成功接收并处理
//3**	重定向，需要进一步的操作以完成请求
//4**	客户端错误，请求包含语法错误或无法完成请求
//5**	服务器错误，服务器在处理请求的过程中发生了错误

// Success 根据官方文档 0 代表成功
var ServerSuccess = newError(0, "Success")

// 200 OK
// var (
// 	ServerSuccess = newError(200, "Success")
// )

// 400 BAD REQUEST
var (
	InvalidParams          = newError(40001, "参数错误")
	PasswordIncorrect      = newError(40002, "密码错误")
	CommentActionUnknown   = newError(40003, "未知的评论操作")
	FavouriteActionUnknown = newError(40004, "未知的收藏操作")
	FileUploadFailed       = newError(40005, "文件上传失败")
	FileIsNotVideo         = newError(40006, "文件不是视频")
	MessageActionUnknown   = newError(40007, "未知的消息操作")
)

// 401 WITHOUT PERMISSION
var (
	NotLogin     = newError(40101, "用户未登录")
	InvalidToken = newError(40102, "无效的Token")
)

// 403 ILLEGAL OPERATION
var (
	PermissionDenied = newError(40301, "权限不足")
	IllegalOperation = newError(40302, "非法操作")
)

// 404 NOT FOUND
var (
	ItemNotFound = newError(40401, "资源不存在")
)

// 409 CONFLICT
var (
	UserAlreadyExist = newError(40901, "用户已存在")
	AlreadyFollowed  = newError(40904, "当前已关注")
	NotFollowed      = newError(40905, "当前未关注")
)

// 500 Server ERROR
var (
	ServerInternal = newError(50001, "服务器内部错误")

	ServerLogic = newError(50301, "逻辑错误")
)

func newError(code int, msg string) ApiErr {
	return ApiErr{
		Code: code,
		Msg:  msg,
	}
}

type ApiErr struct {
	Code int    `json:"status_code"` // 状态码，0-成功，其他值-失败
	Msg  string `json:"status_msg"`  // 返回状态信息
}

// 实现了方法 可以赋值给error
func (e ApiErr) Error() string {
	return e.Msg
}

// WithDetails 在基础错误上追加详细信息，例如：密码错误，密码长度不足6位
func (e ApiErr) WithDetails(detail string) ApiErr {
	return ApiErr{
		Code: e.Code,
		Msg:  e.Msg + ": " + detail,
	}
}
