/*
@Time : 2020/5/18 14:14
@Author : leixianting
@File : yard
@Software: GoLand
*/
package yard

import (
	"chaincode/bsn/constants"
	"chaincode/bsn/enums"
	"chaincode/bsn/errors"
	"chaincode/bsn/handler"
	"chaincode/bsn/request"
	"chaincode/bsn/structs"
	"chaincode/bsn/utils"
	"crypto/md5"
	e "github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 添加新的仓库
func AddNewYard(c router.Context) (interface{}, error) {

	req := c.Param(constants.ROUTER_ADD_NEW_YARD_REQ).(request.AddNewYardReq)

	uid := req.Uid

	address := req.Address

	longitude := req.Longitude

	latitude := req.Latitude

	name := req.Name

	// 生成新的仓库身份地址
	id, err := utils.GenerateAddress()
	if err != nil {
		c.Logger().Errorf("生成地址错误：%s", err.Error())
		return "", e.Wrap(err, errors.ErrPutState)
	}

	now := time.Now()

	var yard = &structs.Yard{
		Id: id,
		Uid: uid,
		Name: name,
		Address: address,
		Longitude: longitude,
		Latitude: latitude,
		CreateTime: now.Unix(),
	}

	// 存储仓库
	err = c.State().Insert(yard)
	if err != nil {
		c.Logger().Errorf("%s：身份地址存储失败", errors.ErrPutState)
		return "", e.Wrap(err, errors.ErrPutState)
	}

	return handler.GenerateSuccessResp(id), nil
}

// 仓库列表
func getYards(uid string, c router.Context) (interface{}, error)  {

	list, err := c.State().List(&structs.Yard{Uid: uid}, structs.Yard{})

	if nil != err {
		c.Logger().Errorf("获取仓库列表错误：%s", err.Error())
		return nil, e.Wrap(err, errors.ErrGetState)
	}

	return list, nil
}

// 入库操作
func AddDonationInbound(c router.Context) (interface{}, error)  {
	req := c.Param(constants.ROUTER_ADD_DONATION_INBOUND_REQ).(request.AddDonationInboundReq)

	donationId := req.DonationId

	result, err := c.State().Get([]string{constants.DONATION_PREFIX, donationId}, &structs.Donation{})
	if nil != err {
		c.Logger().Errorf("获取捐赠单错误：%s", err.Error())
		return nil, e.Wrap(err, errors.ErrGetState)
	}

	donation, ok := result.(structs.Donation)
	if !ok {
		c.Logger().Errorf("获取捐赠单错误")
		return nil, e.Wrap(err, errors.ErrGetState)
	}

	// 查看是否已经入库
	if donation.Status == enums.Finish {
		c.Logger().Errorf("该捐赠单物品已入库")
		return handler.GenerateSuccessResp("该捐赠单物品已入库"), nil
	}

	list := req.List

	sort.Sort(sort.StringSlice(list))
	
	yid := req.Yid

	listEx := donation.List

	sort.Strings(sort.StringSlice(listEx))

	if md5.Sum([]byte(strings.Join(list, ","))) != md5.Sum([]byte(strings.Join(listEx, ","))) {
		c.Logger().Error("物资有出入,捐赠物资:%s,入库物资", listEx, list)
		return handler.GenerateSuccessResp("入库物资与捐赠物资有出入"), nil
	}

	now := time.Now()

	// 物品信息存储
	for _, value := range list {
		c.Logger().Info("清单列表:%s", value)
		result := strings.Split(value, ",")
		var goods = &structs.Goods{
			Yid: yid,
			Lid: donation.ListId,
			CreateTime: now.Unix(),
		}
		for j, v := range result {
			if 0 == j {
				goods.Name = v
			} else if 1 == j {
				goods.Amount, _ = strconv.Atoi(v)
			}
		}

		err := c.State().Insert(goods)
		if err != nil {
			c.Logger().Errorf("存储物资失败：%s", err.Error())
			return nil, e.Wrap(err, errors.ErrPutState)
		}
	}

	donation.Status = enums.Finish

	err = c.State().Put(donation)
	if err != nil {
		c.Logger().Errorf("%s：更新捐赠单子信息失败", donation.Id)
		return "", e.Wrap(err, errors.ErrPutState)
	}
	return handler.GenerateSuccessResp("捐赠物资入库成功"), nil
}

// 出库操作
func AddDonationOutbound(c router.Context) (interface{}, error)  {
	req := c.Param(constants.ROUTER_ADD_DONATION_OUTBOUND_REQ).(request.AddDonationOutboundReq)

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
	if d.Status != enums.Start {
		c.Logger().Info("需求:%s,请不要重复出库", id)
		return handler.GenerateFailResp("请不要重复出库"), nil
	}

	// 检查该需求id的匹配单是否已经存在
	m, err := c.State().Get([]string{constants.MATCH_PREFIX, id}, &structs.Match{})
	if err != nil {
		c.Logger().Error("需求:%s 的匹配单获取失败", id)
		return handler.GenerateFailResp("未找到该需求匹配单"), nil
	}

	match, ok := m.(structs.Match)
	if !ok {
		c.Logger().Error("需求:%s 的匹配单获取失败", id)
		return handler.GenerateFailResp("未找到该需求匹配单"), nil
	}

	list := match.List

	// 出库 清除库存
	for _, value := range list {
		goods := strings.Split(value, ",")
		name := goods[0]
		yid := goods[1]
		lid := goods[3]

		amout, _ := strconv.Atoi(goods[4])

		g, err := c.State().Get([]string{constants.GOODS_PREFIX, name, yid, lid}, &structs.Goods{})
		if err != nil {
			c.Logger().Error("匹配单中货物信息获取失败:%s", id)
			return handler.GenerateFailResp("未找到该需求匹配单"), nil
		}

		gd, ok := g.(structs.Goods)
		if !ok {
			c.Logger().Error("需求:%s 的匹配单获取失败", id)
			return handler.GenerateFailResp("未找到该需求匹配单"), nil
		}

		if gd.Amount - amout == 0 {
			err = c.State().Delete([]string{constants.GOODS_PREFIX, name, yid, lid})
			if err != nil {
				c.Logger().Error("需求:%s 的匹配单获取失败", id)
				return handler.GenerateFailResp("未找到该需求匹配单"), nil
			}
		} else {
			gd.Amount -= amout
			err := c.State().Put(gd)
			if err != nil {
				c.Logger().Error("需求:%s 的匹配单获取失败", id)
				return handler.GenerateFailResp("未找到该需求匹配单"), nil
			}
		}

	}



	d.Status = enums.Processing

	err = c.State().Put(d)
	if err != nil {
		c.Logger().Errorf("%s：需求状态失败", id)
		return "", e.Wrap(err, errors.ErrPutState)
	}

	return handler.GenerateSuccessResp("捐赠物资出库成功"), nil
}

/**
 获取库存
 */
func GetInventory(c router.Context) (interface{}, error) {
	req := c.Param(constants.ROUTER_GET_INVENTORY).(request.GetInventory)

	// 仓库id
	yid := req.Yid

	// 货物名称
	name := req.Name

	if 0 == len(yid) {
		c.Logger().Info("仓库id不能为空")
		return handler.GenerateFailResp("仓库id不能为空"), nil
	}

	if 0 == len(name) {
		c.Logger().Info("货物名称不能为空")
		return handler.GenerateFailResp("货物名称不能为空"), nil
	}

	// 获取到对应的货物
	rs, err := c.State().List([]string{constants.GOODS_PREFIX, name, yid}, &structs.Goods{})

	if err != nil {
		c.Logger().Error("仓库:%s,获取货物:%s失败", yid, name)
		return nil, e.Wrap(err, errors.ErrGetState)
	}

	return rs, nil
}