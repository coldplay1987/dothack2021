/*
@Time : 2020/5/15 15:41
@Author : leixianting
@File : main
@Software: GoLand
*/
package main

import (
	"chaincode/bsn/constants"
	"chaincode/bsn/goods"
	"chaincode/bsn/identity"
	"chaincode/bsn/request"
	"chaincode/bsn/yard"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	_ "github.com/s7techlab/cckit/router/param"
	p "github.com/s7techlab/cckit/router/param"
)


/*
 * 日志打印
 */
var logger = shim.NewLogger("main")

/**
 * 方法路由
 */
func New() *router.Chaincode {
	r := router.New(`carbon_points`)

	r.Init(invokeInit)

	// 添加新用户(捐赠方/受赠方/慈善机构)
	r.Invoke(constants.ROUTER_ADD_NEW_USER_REQ, identity.AddNewUser, p.Struct(constants.ROUTER_ADD_NEW_USER_REQ, &request.AddNewUserReq{}))

	// 添加仓库
	r.Invoke(constants.ROUTER_ADD_NEW_YARD_REQ, yard.AddNewYard, p.Struct(constants.ROUTER_ADD_NEW_YARD_REQ, &request.AddNewYardReq{}))

	// 添加新捐赠
	r.Invoke(constants.ROUTER_ADD_NEW_DONATION_REQ, goods.AddNewDonation, p.Struct(constants.ROUTER_ADD_NEW_DONATION_REQ, &request.AddNewDonationReq{}))

	// 添加新需求
	r.Invoke(constants.ROUTER_ADD_NEW_DEMAND_REQ, goods.AddNewDemand, p.Struct(constants.ROUTER_ADD_NEW_DEMAND_REQ, &request.AddNewDemandReq{}))

	// 获取匹配结果
	r.Invoke(constants.ROUTER_MATERIAL_MATCHING_REQ, goods.MaterialMatching, p.Struct(constants.ROUTER_MATERIAL_MATCHING_REQ, &request.MaterialMatchingReq{}))

	// 物品入库
	r.Invoke(constants.ROUTER_ADD_DONATION_INBOUND_REQ, yard.AddDonationInbound, p.Struct(constants.ROUTER_ADD_DONATION_INBOUND_REQ, &request.AddDonationInboundReq{}))

	// 物品出库
	r.Invoke(constants.ROUTER_ADD_DONATION_OUTBOUND_REQ, yard.AddDonationOutbound, p.Struct(constants.ROUTER_ADD_DONATION_OUTBOUND_REQ, &request.AddDonationOutboundReq{}))

	// 物品签收
	r.Invoke(constants.ROUTER_ADD_RECEIPT_REQ, goods.Receipt, p.Struct(constants.ROUTER_ADD_RECEIPT_REQ, &request.AddReceiptReq{}))

	// 物品库存
	r.Invoke(constants.ROUTER_GET_INVENTORY, yard.GetInventory, p.Struct(constants.ROUTER_GET_INVENTORY, &request.GetInventory{}))

	return router.NewChaincode(r)
}

/**
 * chaincode的init方法
 * 设置了一个合约所有者
 */
func invokeInit(c router.Context) (interface{}, error) {
	return owner.SetFromCreator(c)
}

/**
 * main方法
 */
func main() {
	cc := New()
	if err := shim.Start(cc); err != nil {
		logger.Error("合约启动失败: %s", err)
	}
}

