/*
@Time : 2020/5/17 15:54
@Author : leixianting
@File : goods
@Software: GoLand
*/
package goods

import (
	"bytes"
	"chaincode/bsn/constants"
	"chaincode/bsn/enums"
	"chaincode/bsn/errors"
	"chaincode/bsn/handler"
	"chaincode/bsn/request"
	"chaincode/bsn/structs"
	"chaincode/bsn/utils"
	e "github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
	"strconv"
	"strings"
	"time"
)

/**
添加新的捐赠
*/
func AddNewDonation(c router.Context) (interface{}, error) {
	req := c.Param(constants.ROUTER_ADD_NEW_DONATION_REQ).(request.AddNewDonationReq)

	uid := req.Uid

	listId := req.ListId

	list := req.List

	address := req.Address

	longitude := req.Longitude

	latitude := req.Latitude

	now := time.Now()

	// 生成新id
	id, err := utils.GenerateAddress()
	if err != nil {
		c.Logger().Errorf("生成捐赠地址错误：%s", err.Error())
		return "", e.Wrap(err, errors.ErrPutState)
	}

	var donation = &structs.Donation{
		Id: id,

		Uid: uid,

		ListId: listId,

		List: list,

		Address: address,

		Longitude: longitude,

		Latitude: latitude,

		Status: enums.Start,

		CreateTime: now.Unix(),
	}

	err = c.State().Insert(donation)
	if err != nil {
		c.Logger().Errorf("%s：捐赠清单存储失败", id)
		return "", e.Wrap(err, errors.ErrPutState)
	}

	return handler.GenerateSuccessResp(id), nil
}

/**
添加新的需求
*/
func AddNewDemand(c router.Context) (interface{}, error) {

	req := c.Param(constants.ROUTER_ADD_NEW_DEMAND_REQ).(request.AddNewDemandReq)

	uid := req.Uid

	listId := req.ListId

	list := req.List

	address := req.Address

	longitude := req.Longitude

	latitude := req.Latitude

	charityId := req.CharityId

	now := time.Now()

	// 生成新id
	id, err := utils.GenerateAddress()
	if err != nil {
		c.Logger().Errorf("生成需求id错误：%s", err.Error())
		return "", e.Wrap(err, errors.ErrPutState)
	}

	var demand = &structs.Demand {
		Id: id,

		Uid: uid,

		CharityId: charityId,

		ListId: listId,

		List: list,

		Address: address,

		Longitude: longitude,

		Latitude: latitude,

		CreateTime: now.Unix(),
	}

	err = c.State().Insert(demand)
	if err != nil {
		c.Logger().Errorf("%s：需求清单存储失败", id)
		return "", e.Wrap(err, errors.ErrPutState)
	}

	return handler.GenerateSuccessResp(id), nil
}

// 物资匹配
func MaterialMatching(c router.Context) (interface{}, error) {
	req := c.Param(constants.ROUTER_MATERIAL_MATCHING_REQ).(request.MaterialMatchingReq)

	// 需求id
	id := req.Id

	// 慈善机构id
	charityId := req.CharityId

	// 检查该需求id的匹配单是否已经存在
	result, err := c.State().Get([]string{constants.MATCH_PREFIX, id}, &structs.Match{})

	if result != nil {
		 c.Logger().Info("该需求的匹配单已经存在,直接返回")
		 return handler.GenerateSuccessResp(result), nil
	}

	// 检查该需求id是否已经存在
	result, err = c.State().Get([]string{constants.DEMAND_PREFIX, id}, &structs.Demand{})

	if err != nil {
		c.Logger().Errorf("获取需求单失败：%s", err.Error())
		return "", e.Wrap(err, errors.ErrGetState)
	}

	demand, ok := result.(structs.Demand)

	if !ok {
		c.Logger().Errorf("获取需求单失败")
		return "", e.Wrap(err, errors.ErrGetState)
	}

	list := demand.List


	// 先获取到该慈善机构所有的仓库地址
	rs, err := c.State().List([]string{constants.YARD_PREFIX, charityId}, &structs.Yard{})
	yList, ok := rs.([]interface{})
	if !ok {
		c.Logger().Error("匹配失败")
		return handler.GenerateFailResp("匹配失败"), nil
	}

	var yards = make([]string, 0)

	for _, yd := range yList{
		yard, ok := yd.(structs.Yard)
		if !ok {
			c.Logger().Error("匹配失败")
			return handler.GenerateFailResp("匹配失败"), nil
		}
		la1, _ := strconv.ParseFloat(yard.Latitude, 64)
		lgt1, _ := strconv.ParseFloat(yard.Longitude, 64)
		la2, _ := strconv.ParseFloat(demand.Latitude, 64)
		lgt2, _ := strconv.ParseFloat(demand.Longitude, 64)
		distance := utils.LatitudeLongitudeDistance(la1, lgt1, la2, lgt2)
		d := strconv.FormatFloat(distance, 'f', -1, 64)
		var buffer bytes.Buffer
		buffer.WriteString(d)
		buffer.WriteString(",")
		buffer.WriteString(yard.Id)
		buffer.WriteString(",")
		buffer.WriteString(yard.Name)
		yards = append(yards, buffer.String())
	}

	// 对距离排序
	var temp string
	for i := 0;i < len(yards) - 1;i++{
		for j := 0 ;j<len(yards) -1 - i;j++{

			y1 := strings.Split(yards[j], ",")
			y2 := strings.Split(yards[j+1], ",")
			d1,err  := strconv.ParseFloat(y1[0], 64)
			if err != nil {
				c.Logger().Error("匹配失败")
				return handler.GenerateFailResp("匹配失败"), nil
			}
			d2,err  := strconv.ParseFloat(y2[0], 64)
			if err != nil {
				c.Logger().Error("匹配失败")
				return handler.GenerateFailResp("匹配失败"), nil
			}

			if d1 > d2 {
				temp = yards[j]
				yards[j] = yards[j+1]
				yards[j+1] = temp
			}
		}
	}

	var matchList = make([]string, 0)

	// 需求物资与仓库存储物资进行匹配
	for _, value := range list {

		result := strings.Split(value, ",")

		name := result[0]
		amountNeed, _ := strconv.Atoi(result[1])

		total := 0

		// 匹配物资
		for _, yd := range yards {
			y := strings.Split(yd, ",")
			yid := y[1]
			yName := y[2]
			// 获取到对应的货物
			rs, err := c.State().List([]string{constants.GOODS_PREFIX, name, yid}, &structs.Goods{})
			if err != nil {
				c.Logger().Error("匹配失败")
				return handler.GenerateFailResp("匹配失败"), nil
			}
			goodsList, ok := rs.([]interface{})
			if !ok {
				c.Logger().Error("匹配失败")
				return handler.GenerateFailResp("匹配失败"), nil
			}

			for _, va := range goodsList {

				var buffer bytes.Buffer

				g, ok := va.(structs.Goods)
				if !ok {
					c.Logger().Error("匹配失败")
					return handler.GenerateFailResp("匹配失败"), nil
				}

				if total >= amountNeed {
					c.Logger().Info("%s:匹配结束", g.Name)
					break
				}

				amount := g.Amount

				buffer.WriteString(g.Name)
				buffer.WriteString(",")
				buffer.WriteString(yid)
				buffer.WriteString(",")
				buffer.WriteString(yName)
				buffer.WriteString(",")
				buffer.WriteString(g.Lid)
				buffer.WriteString(",")

				var lockAmount = 0
				// 小于需求量 全部需要
				if amount < amountNeed - total {
					buffer.WriteString(strconv.Itoa(amount))
					lockAmount = amount
				} else {
					buffer.WriteString(strconv.Itoa(amountNeed))
					lockAmount = amountNeed
				}
				matchList = append(matchList, buffer.String())
				total += lockAmount

				g.Amount -= lockAmount

				g.Lock = append(g.Lock, id + "," + strconv.Itoa(lockAmount))

				err := c.State().Put(g)
				if err != nil {
					c.Logger().Errorf("%s：匹配结果存储失败", id)
					return "", e.Wrap(err, errors.ErrPutState)
				}
			}
		}

	}

	var match = &structs.Match{
		Did: id,
		CharityId: charityId,
	}

	match.List = matchList

	match.CreateTime = time.Now().Unix()

	err = c.State().Insert(match)
	if err != nil {
		c.Logger().Errorf("%s：匹配结果存储失败", id)
		return "", e.Wrap(err, errors.ErrPutState)
	}

	return handler.GenerateSuccessResp(match), nil
}


// 物资签收
func Receipt(c router.Context) (interface{}, error) {
	req := c.Param(constants.ROUTER_ADD_RECEIPT_REQ).(request.AddReceiptReq)

	// 需求id
	id := req.Id

	// 检查该需求id是否已经存在
	demand, err := c.State().Get([]string{constants.DEMAND_PREFIX, id}, &structs.Demand{})
	if err != nil {
		c.Logger().Errorf("获取需求单失败：%s", err.Error())
		return "", e.Wrap(err, errors.ErrGetState)
	}
	d, ok := demand.(structs.Demand)
	if !ok {
		c.Logger().Error("需求:%s 的匹配单获取失败", id)
		return handler.GenerateFailResp("未找到该需求匹配单"), nil
	}

	d.Status = enums.Finish

	err = c.State().Put(d)
	if err != nil {
		c.Logger().Errorf("%s：需求状态失败", id)
		return "", e.Wrap(err, errors.ErrPutState)
	}

	return handler.GenerateSuccessResp("签收成功"), nil

}



