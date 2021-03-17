/*
@Time : 2020/5/15 16:30
@Author : leixianting
@File : struct
@Software: GoLand
*/
package structs

import (
	"chaincode/bsn/constants"
	"chaincode/bsn/enums"
)

type User struct {

	/*
		用户id
	*/
	Uid string `json:"uid"`

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

	/**
	创建时间
	 */
	CreateTime int64 `json:"createTime"`
}

// 存储账户的key
func (user User) Key() ([]string, error) {
	return []string{constants.USER_PREFIX, user.Uid}, nil
}

// 捐赠
type Donation struct {

	Id string `json:"id"`
	
	// 用户链上唯一id
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
	
	// 状态
	Status enums.Status `json:"status"`

	/**
	创建时间
	*/
	CreateTime int64 `json:"createTime"`
}

// 存储捐赠的key
func (d Donation) Key() ([]string, error) {
	return []string{constants.DONATION_PREFIX, d.Id}, nil
}

// 需求
type Demand struct {
	Id string `json:"id"`

	// 用户链上唯一地址
	Uid string `json:"uid"`

	// 慈善机构id
	CharityId string `json:"charityId"`

	// 物资清单id
	ListId string `json:"listId"`

	// 物资列表[ "物资名称,  物资数量", "物资名称, 物资数量"]
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

	// 状态
	Status enums.Status `json:"status"`

	/**
	创建时间
	*/
	CreateTime int64 `json:"createTime"`
}

// 存储需求的key
func (d Demand) Key() ([]string, error) {
	return []string{constants.DEMAND_PREFIX, d.Id}, nil
}

// 货物
type Goods struct {
	Id string `json:"id"`
	
	Name string `json:"name"`
	
	Amount int `json:"amount"`

	// 匹配锁定列表
	// ["需求id,数量", "需求id,数量"]
	Lock []string `json:"lock"`

	// 仓库id
	Yid string `json:"yid"`

	// 清单id
	Lid string `json:"lid"`

	/**
	创建时间
	*/
	CreateTime int64 `json:"createTime"`
}

// 存储货物的key
func (g Goods) Key() ([]string, error) {
	return []string{constants.GOODS_PREFIX, g.Name, g.Yid, g.Lid}, nil
}

// 匹配结果
type Match struct {
	// 需求id
	Did string `json:"did"`

	// 慈善机构id
	CharityId string `json:"charityId"`

	// 匹配结果
	List []string `json:"list"`

	/**
	创建时间
	*/
	CreateTime int64 `json:"createTime"`
}

// 存储货物的key
func (m Match) Key() ([]string, error) {
	return []string{constants.MATCH_PREFIX, m.Did}, nil
}

// 仓库
type Yard struct {

	Id string `json:"id"`

	// 所属人的账户地址
	Uid string `json:"uid"`

	// 仓库
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

	/**
	创建时间
	*/
	CreateTime int64 `json:"createTime"`
}

// 存储仓库的key
func (y Yard) Key() ([]string, error) {
	return []string{constants.YARD_PREFIX, y.Uid, y.Id}, nil
}