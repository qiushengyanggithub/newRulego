package handlers

import (
	"github.com/rulego/rulego/api/types"
	"strings"
)

type HelloWord struct{}

func (n *HelloWord) Type() string {
	return "test/hello"
}
func (n *HelloWord) New() types.Node {
	return &HelloWord{}
}

// OnMsg 处理消息
func (n *HelloWord) Init(ruleConfig types.Config, configuration types.Configuration) error {

	return nil
}

// OnMsg 处理消息
// OnMsg 是 HelloWord 类型的方法，用于处理接收到的消息。
// 该方法将消息数据进行转换，并发送成功通知。
// 参数:
//
//	ctx types.RuleContext: 上下文环境，提供了访问规则引擎功能的接口。
//	msg types.RuleMsg: 接收到的消息，包含待处理的数据。
func (n *HelloWord) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	// 获取原始消息数据
	data := msg.GetData()
	// 将数据转换为大写，并在前面添加"HelloWord"
	modifiedData := "HelloWord" + strings.ToUpper(data)
	// 修改消息内容
	msg.SetData(modifiedData)
	// 通知上下文环境，消息处理成功
	ctx.TellSuccess(msg)
}

func (n *HelloWord) Destroy() {
	// 释放资源
}
