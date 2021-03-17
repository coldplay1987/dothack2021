/*
@Time : 2020/5/15 15:52
@Author : leixianting
@File : request
@Software: GoLand
*/
package request

type AddNewUserReq struct {

	/*
	用户类型
	 */
	UType string `json:"uType"`

	/**
	用户名
	 */
	UName string `json:"uName"`

	/**
	联系人
	 */
	Contact string `json:"contact"`

	/**
	手机号码
	 */
	Telephone string `json:"telephone"`

	/**
	经度
	 */
	Longitude string `json:"longitude"`

	/**
	纬度
	 */
	Latitude string `json:"latitude"`
}

// 捐赠
type AddNewDonationReq struct {

	// 用户链上唯一地址
	Uid string `json:"uid"`

	// 物资清单id
	ListId string `json:"listId"`

	// 物资列表,格式:[ "物资名称, 物资数量", "物资名称, 物资数量"]
	List []string `json:"list"`

	// 从哪里运来的
	// 浙江-杭州-萧山
	Address string `json:"address"`

	/**
	经度
	*/
	Longitude string `json:"longitude"`

	/**
	纬度
	*/
	Latitude string `json:"latitude"`
}

// 入库
type AddDonationInboundReq struct {

	// 捐赠id
	DonationId string `json:"donationId"`
	
	// 仓库id
	Yid string `json:"yid"`

	// 物资列表,格式:[ "物资名称, 物资数量", "物资名称, 物资数量"]
	List []string `json:"list"`
	
}

// 出库请求
type AddDonationOutboundReq struct {
	// 需求id
	Id string `json:"id"`
}

// 需求
type AddNewDemandReq struct {

	// 用户链上唯一地址
	Uid string `json:"uid"`

	// 慈善机构id
	CharityId string `json:"charityId"`

	// 物资清单id
	ListId string `json:"listId"`

	// 物资列表,[ "物资名称, 物资数量", "物资名称, 物资数量"]
	List []string `json:"list"`

	// 运往哪里
	// 浙江-杭州-萧山
	Address string `json:"address"`

	/**
	经度
	*/
	Longitude string `json:"longitude"`

	/**
	纬度
	*/
	Latitude string `json:"latitude"`
}

type MaterialMatchingReq struct {
	// 需求id
	Id string `json:"id"`
	
	// 慈善机构id
	CharityId string `json:"charityId"`
}

// 添加新的仓库请求
type AddNewYardReq struct {
	// 所属人的账户id
	Uid string `json:"uid"`

	// 仓库名称
	Name string `json:"name"`

	// 仓库地址
	Address string `json:"address"`

	/**
	经度
	*/
	Longitude string `json:"longitude"`

	/**
	纬度
	*/
	Latitude string `json:"latitude"`
}

// 签收
type AddReceiptReq struct {
	// 需求id
	Id string `json:"id"`
}

type GetInventory struct {
	// 仓库id
	Yid string `json:"yid"`
	
	// 物品名称
	Name string `json:"name"`
}