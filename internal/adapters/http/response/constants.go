// Package response 定义响应消息常量
package response

const (
	// MsgSuccess 表示操作成功
	MsgSuccess = "操作成功"

	// MsgCreated 表示资源创建成功
	MsgCreated = "创建成功"

	// MsgUpdated 表示资源更新成功
	MsgUpdated = "更新成功"

	// MsgDeleted 表示资源删除成功
	MsgDeleted = "删除成功"

	// MsgValidationFailed 表示请求参数验证失败
	MsgValidationFailed = "验证失败"

	// MsgAuthenticationRequired 表示需要身份认证
	MsgAuthenticationRequired = "请先登录"

	// MsgAccessForbidden 表示无权限访问
	MsgAccessForbidden = "无权访问"

	// MsgResourceNotFound 表示请求的资源不存在
	MsgResourceNotFound = "资源不存在"

	// MsgResourceConflict 表示资源冲突
	MsgResourceConflict = "资源冲突"

	// MsgInternalError 表示服务器内部错误
	MsgInternalError = "服务器内部错误"

	// MsgServiceUnavailable 表示服务暂时不可用
	MsgServiceUnavailable = "服务暂时不可用"

	// MsgRateLimitExceeded 表示请求过于频繁
	MsgRateLimitExceeded = "请求过于频繁"
)
