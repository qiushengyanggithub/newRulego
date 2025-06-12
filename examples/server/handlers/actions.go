// handlers/handlers.go

package handlers

import (
	"github.com/rulego/rulego/api/types"
	"strings"
)

type UpperNode struct{}

func (n *UpperNode) Type() string {
	return "test/upper"
}

func (n *UpperNode) New() types.Node {
	return &UpperNode{}
}

func (n *UpperNode) Init(ruleConfig types.Config, configuration types.Configuration) error {

	return nil
}

func (n *UpperNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	// 获取原始数据并转换为大写
	data := msg.GetData()
	modifiedData := strings.ToUpper(data)

	// 修改消息内容
	msg.SetData(modifiedData)

	// 将修改后的消息发送到下一个节点
	ctx.TellSuccess(msg)
}

func (n *UpperNode) Destroy() {
	// 释放资源
}
