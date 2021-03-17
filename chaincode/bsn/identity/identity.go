/*
@Time : 2020/5/15 15:51
@Author : leixianting
@File : identity
@Software: GoLand
*/
package identity

import (
	"chaincode/bsn/constants"
	"chaincode/bsn/errors"
	"chaincode/bsn/handler"
	"chaincode/bsn/request"
	"chaincode/bsn/structs"
	"chaincode/bsn/utils"
	e "github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
	"time"
)

/**
添加用户,生成链上唯一id
*/
func AddNewUser(c router.Context) (interface{}, error) {

	req := c.Param(constants.ROUTER_ADD_NEW_USER_REQ).(request.AddNewUserReq)

	// 用户名
	uName := req.UName

	// 用户类型
	uType := req.UType

	// 用户手机
	telephone := req.Telephone

	// 联系热
	contact := req.Contact

	// 经度
	longitude := req.Longitude

	// 纬度
	latitude := req.Latitude

	// 生成新的身份地址
	uid, err := utils.GenerateAddress()
	if err != nil {
		c.Logger().Errorf("生成地址错误：%s", err.Error())
		return "", e.Wrap(err, errors.ErrPutState)
	}

	now := time.Now()
	var user = &structs.User{
		Uid:        uid,
		UName:      uName,
		UType:      uType,
		Telephone:  telephone,
		Contact:    contact,
		Longitude:  longitude,
		Latitude:   latitude,
		CreateTime: now.Unix(),
	}

	// 存储用户
	err = c.State().Insert(user)
	if err != nil {
		c.Logger().Errorf("%s：身份地址存储失败", errors.ErrPutState)
		return "", e.Wrap(err, errors.ErrPutState)
	}

	return handler.GenerateSuccessResp(uid), nil
}


