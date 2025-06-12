package handlers

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/robinson/gos7"
	"github.com/rulego/rulego/api/types"
	"math"
	"strconv"
	"time"
)

// S7_SMART200_Read 从 S7-200 SMART PLC 中读取数据。
// 该函数根据提供的 JSON 数据解析出 IP 地址、起始地址、读取数量和数据类型，
// 然后连接到 PLC 并读取相应数量的数据。
func S7_SMART200_Read(ctx types.RuleContext, msg types.RuleMsg) {
	messageData := msg.GetData()

	// 解析 JSON 数据
	data, err := parseJSON(messageData)
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("解析 JSON 数据时出错: %v", err), nil)
		return
	}

	// 获取并验证 IP 地址
	ipStr, ok := data["ip"].(string)
	if !ok {
		handleError(ctx, msg, fmt.Errorf("ip 值类型不正确"), nil)
		return
	}

	// 获取起始地址
	start, err := getInteger(data, "start")
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("获取 start 值时出错: %v", err), nil)
		return
	}

	// 获取读取字节长度
	count, err := getInteger(data, "count")
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("获取 count 值时出错: %v", err), nil)
		return
	}

	// 获取并验证数据类型
	typesStr, ok := data["types"].(string)
	if !ok {
		handleError(ctx, msg, fmt.Errorf("types 值类型不正确"), nil)
		return
	}

	// 创建 TCP 客户端处理器并设置相关参数
	handler := gos7.NewTCPClientHandler(ipStr, 0, 1)
	handler.Timeout = 5 * time.Second
	handler.IdleTimeout = 10 * time.Second
	handler.PDULength = 960

	// 连接到 PLC
	err = handler.Connect()
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("连接失败: %v", err), nil)
		return
	}
	defer func() {
		if err := handler.Close(); err != nil {
			fmt.Printf("关闭连接时出错: %v", err)
		}
	}()

	// 创建客户端实例
	client := gos7.NewClient(handler)

	// 读取数据
	VBuffer := make([]byte, count)
	MBuffer := make([]byte, count)
	IBuffer := make([]byte, count)
	QBuffer := make([]byte, count)

	err = client.AGReadDB(1, start, count, VBuffer)
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("读取错误: %v", err), nil)
		return
	}

	// 根据寄存器的数据类型处理读取的结果
	var result interface{}
	switch typesStr {
	case "VB":
		result = VBuffer[0]
	case "VW":
		result = binary.BigEndian.Uint16(VBuffer[0:2])
	case "VD":
		result = binary.BigEndian.Uint32(VBuffer[0:4])
	case "VF":
		result = math.Float32frombits(binary.BigEndian.Uint32(VBuffer[0:4]))
	case "MB":
		result = MBuffer[0]
	case "MW":
		result = binary.BigEndian.Uint16(MBuffer[0:2])
	case "MD":
		result = binary.BigEndian.Uint32(MBuffer[0:4])
	case "MF":
		result = math.Float32frombits(binary.BigEndian.Uint32(MBuffer[0:4]))
	case "IB":
		result = IBuffer[0]
	case "QB":
		result = QBuffer[0]
	default:
		handleError(ctx, msg, fmt.Errorf("寄存器类型未知类型: %s", typesStr), nil)
		return
	}

	// 将结果转换为 JSON 格式并更新消息数据
	resultMap := map[string]interface{}{"value": result}
	msgDataBytes, err := json.Marshal(resultMap)
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("转换 JSON 时出错: %v", err), result)
		return
	}

	// ✅ 使用 SetData 替代直接赋值
	msg.SetData(string(msgDataBytes))

	ctx.TellSuccess(msg)
}

// S7_SMART200_Write 向 S7-200 SMART PLC 写入数据。
// 该函数解析 JSON 数据以获取 IP 地址、数据类型、写入数量、值和起始地址，
// 然后连接到 PLC 并写入相应数量的数据。
func S7_SMART200_Write(ctx types.RuleContext, msg types.RuleMsg) {
	messageData := msg.GetData()

	// 解析 JSON 数据
	data, err := parseJSON(messageData)
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("解析 JSON 数据时出错: %v", err), nil)
		return
	}

	// 获取并验证 IP 地址
	ipStr, ok := data["ip"].(string)
	if !ok {
		handleError(ctx, msg, fmt.Errorf("ip 值类型不正确"), nil)
		return
	}

	// 获取并验证数据类型
	typesStr, ok := data["types"].(string)
	if !ok {
		handleError(ctx, msg, fmt.Errorf("types 值类型不正确"), nil)
		return
	}

	// 获取写入数量
	count, err := getInteger(data, "count")
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("获取 count 值时出错: %v", err), nil)
		return
	}

	// 获取写入的值
	valueStr, ok := data["value"].(string)
	if !ok {
		handleError(ctx, msg, fmt.Errorf("value 值类型不正确"), nil)
		return
	}

	// 获取起始地址
	start, err := getInteger(data, "start")
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("获取 start 值时出错: %v", err), nil)
		return
	}

	// 准备写入的数据缓冲区
	VBuffer := make([]byte, count)
	MBuffer := make([]byte, count)
	QBuffer := make([]byte, count)

	switch typesStr {
	case "VB":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("转换 valueStr 为整数时出错: %v", err), nil)
			return
		}
		VBuffer[0] = byte(value)
	case "VW":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("转换 valueStr 为整数时出错: %v", err), nil)
			return
		}
		binary.BigEndian.PutUint16(VBuffer[0:2], uint16(value))
	case "VD":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("转换 valueStr 为整数时出错: %v", err), nil)
			return
		}
		binary.BigEndian.PutUint32(VBuffer[0:4], uint32(value))
	case "VF":
		value, err := strconv.ParseFloat(valueStr, 32)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("转换 valueStr 为浮点数时出错: %v", err), nil)
			return
		}
		binary.BigEndian.PutUint32(VBuffer[0:4], math.Float32bits(float32(value)))
	case "MB":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("转换 valueStr 为整数时出错: %v", err), nil)
			return
		}
		MBuffer[0] = byte(value)
	case "MW":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("转换 valueStr 为整数时出错: %v", err), nil)
			return
		}
		binary.BigEndian.PutUint16(MBuffer[0:2], uint16(value))
	case "MD":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("转换 valueStr 为整数时出错: %v", err), nil)
			return
		}
		binary.BigEndian.PutUint32(MBuffer[0:4], uint32(value))
	case "MF":
		value, err := strconv.ParseFloat(valueStr, 32)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("转换 valueStr 为浮点数时出错: %v", err), nil)
			return
		}
		binary.BigEndian.PutUint32(MBuffer[0:4], math.Float32bits(float32(value)))
	case "QB":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("转换 valueStr 为整数时出错: %v", err), nil)
			return
		}
		QBuffer[0] = byte(value)
	default:
		handleError(ctx, msg, fmt.Errorf("寄存器类型未知类型: %s", typesStr), nil)
		return
	}

	// 创建 TCP 客户端处理器并设置相关参数
	handler := gos7.NewTCPClientHandler(ipStr, 0, 1)
	handler.Timeout = 5 * time.Second
	handler.IdleTimeout = 10 * time.Second
	handler.PDULength = 960

	// 连接到 PLC
	err = handler.Connect()
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("连接失败: %v", err), nil)
		return
	}
	defer func() {
		if err := handler.Close(); err != nil {
			fmt.Printf("关闭连接时出错: %v", err)
		}
	}()

	// 创建客户端实例
	client := gos7.NewClient(handler)

	// 写入数据
	err = client.AGWriteDB(1, start, count, VBuffer)
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("写入错误: %v", err), nil)
		return
	}

	// 更新消息数据并通知成功
	result := map[string]interface{}{"value": valueStr}
	msgDataBytes, err := json.Marshal(result)
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("转换 JSON 时出错: %v", err), result)
		return
	}

	// ✅ 使用 SetData 替代直接赋值
	msg.SetData(string(msgDataBytes))

	ctx.TellSuccess(msg)
}

// handleError 处理错误并推送到下游节点
func handleError(ctx types.RuleContext, msg types.RuleMsg, err error, result interface{}) {
	resultsMap := map[string]interface{}{
		"error": err.Error(),
		"value": result,
	}
	msgDataBytes, _ := json.Marshal(resultsMap)
	msg.SetData(string(msgDataBytes))
	ctx.TellFailure(msg, err)
}

// parseJSON 解析 JSON 字符串并返回解析后的数据。
func parseJSON(data string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 数据时出错: %v", err)
	}
	return result, nil
}

// getInteger 从数据映射中获取指定键的整数值。
func getInteger(data map[string]interface{}, key string) (int, error) {
	value, ok := data[key].(string)
	if !ok {
		return 0, fmt.Errorf("%s 值类型不正确", key)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("转换 %s 为整数时出错: %v", key, err)
	}
	return intValue, nil
}
