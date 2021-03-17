/*
@Time : 2020/5/15 16:40
@Author : leixianting
@File : error
@Software: GoLand
*/
package errors

const (
	ErrInvalidArgs        = "参数错误"   //errors.New("参数错误")
	ErrInvalidFunction    = "方法名错误"  //errors.New("方法名错误")
	ErrUpdate             = "信息更新失败" //errors.New("信息更新失败")
	ErrPutState           = "信息存储失败"
	ErrGetState           = "信息获取失败"
	ErrMarshal            = "序列化失败"
	ErrUnMarshal          = "反序列化失败"
	ErrCreateCompositeKey = "获取复合键失败"
)
