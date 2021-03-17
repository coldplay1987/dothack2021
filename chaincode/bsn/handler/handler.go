/*
@Time : 2020/5/18 16:58
@Author : leixianting
@File : handler
@Software: GoLand
*/
package handler

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func GenerateSuccessResp(data interface{}) Response {
	resp := Response{
		Code:    0,
		Message: "操作成功",
		Data:    data,
	}

	return resp
}

func GenerateFailResp(data interface{}) Response {
	resp := Response{
		Code:    -200,
		Message: "操作失败",
		Data:    data,
	}

	return resp
}