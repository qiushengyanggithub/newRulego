package handlers

import (
	"fmt"
	"github.com/rulego/rulego/api/types"
)

type ToStringNode struct{}

func (n *ToStringNode) Type() string {
	return "test/toString"
}

func (n *ToStringNode) New() types.Node {
	return &ToStringNode{}
}

func (n *ToStringNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	return nil
}

func (n *ToStringNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	// 安全地将 msg.Data 转换为字符串
	msg.SetData(fmt.Sprintf("%v", msg.GetData()))
	ctx.TellSuccess(msg)
}

func (n *ToStringNode) Destroy() {
	// 释放资源
}
