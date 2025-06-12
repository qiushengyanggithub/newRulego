package handlers

import (
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/components/action"
)

func RegisterCustomComponentsAndFunctions() {
	// 注册自定义函数
	action.Functions.Register("Modbus-Read", ModbusTCP_Read)       // 注册自定义函数读取单寄存器
	action.Functions.Register("Modbus-Write", ModbusTCP_Write)     // 注册自定义函数写入单寄存器
	action.Functions.Register("SMART200-Read", S7_SMART200_Read)   // 注册自定义函数读取S7-SMART200数据
	action.Functions.Register("SMART200-Write", S7_SMART200_Write) // 注册自定义函数写入S7-SMART200数据

	// 注册自定义组件
	rulego.Registry.Register(&UpperNode{})    // 注册自定义组件大写转换器
	rulego.Registry.Register(&ToStringNode{}) // 注册自定义组件添加字符串信息
	rulego.Registry.Register(&HelloWord{})    // 注册自定义组件添加HelloWorld信息
}
