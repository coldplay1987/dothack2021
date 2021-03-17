/*
@Time : 2020/5/18 17:47
@Author : leixianting
@File : enum
@Software: GoLand
*/
package enums

type Status uint32

const (
	Start Status = iota
	Processing
	Finish
)
