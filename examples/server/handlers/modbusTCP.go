package handlers

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/goburrow/modbus"
	"github.com/shopspring/decimal"
	"math"
	"strconv"
	"time"
	"unicode"

	"github.com/rulego/rulego/api/types"
)

// extractField 提取单个字段的公共逻辑
func extractField(key string, messageData string) (string, error) {
	data, err := TransformData(key, []byte(messageData))
	if err != nil {
		return "", err
	}
	value, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("%s 值类型不正确", key)
	}
	return value, nil
}

// extractCommonFields 提取消息数据中的公共字段
func extractCommonFields(messageData string) (ip string, port string, slaveid byte, typess string, start uint16, order string, value string, extract uint16, err error) {
	ip, err = extractField("ip", messageData)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}

	port, err = extractField("port", messageData)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}

	idStr, err := extractField("slaveid", messageData)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}
	slaveid = byte(idInt)

	typess, err = extractField("types", messageData)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}

	startStr, err := extractField("start", messageData)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}
	startInt, err := strconv.Atoi(startStr)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}
	start = uint16(startInt)

	order, err = extractField("order", messageData)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}

	value, err = extractField("value", messageData)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}

	extractStr, err := extractField("extract", messageData)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}
	extractInt, err := strconv.Atoi(extractStr)
	if err != nil {
		return "", "", 0, "", 0, "", "", 0, err
	}
	extract = uint16(extractInt)

	return ip, port, slaveid, typess, start, order, value, extract, nil
}

// shortsToInt16Array 将读取的结果转换为 int16 数组
func shortsToInt16Array(results []byte) []int16 {
	result := make([]int16, len(results)/2)
	for i := 0; i < len(results)/2; i++ {
		result[i] = int16(uint16(results[i*2])<<8 | uint16(results[i*2+1]))
	}
	return result
}

// checkResultLength 检查结果长度是否符合预期
func checkResultLength(results []byte, expectedLength int) error {
	if len(results) < expectedLength {
		return fmt.Errorf("结果长度不足")
	}
	return nil
}

// ModbusTCP_Read 函数用于读取 Modbus TCP 协议的数据
func ModbusTCP_Read(ctx types.RuleContext, msg types.RuleMsg) {
	ip, port, id, typess, start, order, _, extract, err := extractCommonFields(msg.GetData())
	if err != nil {
		handleError(ctx, msg, err, nil)
		return
	}

	handler := modbus.NewTCPClientHandler(ip + ":" + port)
	handler.Timeout = 10 * time.Second
	handler.SlaveId = id

	err = handler.Connect()
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("Connect failed: %s", err.Error()), nil)
		return
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	var results []byte
	var result interface{}
	var readErr error

	firstDigit := ""
	for _, char := range typess {
		if unicode.IsDigit(char) {
			firstDigit = string(char)
			break
		}
	}

	switch firstDigit {
	case "0":
		results, err = client.ReadCoils(start, 1)
	case "1":
		results, err = client.ReadDiscreteInputs(start, 1)
	case "3":
		switch typess {
		case "3xbyte", "3xint", "3xuint", "3xbit":
			results, err = client.ReadInputRegisters(start, 1)
		case "3xDint", "3xfloat":
			results, err = client.ReadInputRegisters(start, 2)
		default:
			handleError(ctx, msg, fmt.Errorf("不支持的类型"), nil)
			return
		}
	case "4":
		switch typess {
		case "4xbyte", "4xint", "4xuint", "4xbit":
			results, err = client.ReadHoldingRegisters(start, 1)
		case "4xDint", "4xfloat":
			results, err = client.ReadHoldingRegisters(start, 2)
		default:
			handleError(ctx, msg, fmt.Errorf("不支持的类型"), nil)
			return
		}
	}

	if err != nil {
		handleError(ctx, msg, fmt.Errorf("读取寄存器失败: %s", err.Error()), nil)
		return
	}

	switch typess {
	case "0xbit", "1xbit":
		if err := checkResultLength(results, 1); err != nil {
			readErr = err
		} else {
			result = fmt.Sprintf("%d", results[0])
		}
	case "3xbyte", "4xbyte":
		if err := checkResultLength(results, 1); err != nil {
			readErr = err
		} else {
			result = fmt.Sprintf("%d", results[0])
		}
	case "3xint", "4xint":
		if err := checkResultLength(results, 2); err != nil {
			readErr = err
		} else {
			adjustedResults := make([]byte, len(results))
			for i := 0; i < len(results); i += 2 {
				adjustedResults[i] = results[i+1]
				adjustedResults[i+1] = results[i]
			}
			err = reverseBytes2(adjustedResults, order)
			if err != nil {
				readErr = err
			} else {
				result = fmt.Sprintf("%d", shortsToInt16Array(adjustedResults)[0])
			}
		}
	case "3xuint", "4xuint":
		if err := checkResultLength(results, 2); err != nil {
			readErr = err
		} else {
			err := reverseBytes2(results, order)
			if err != nil {
				readErr = err
			} else {
				result = fmt.Sprintf("%d", binary.BigEndian.Uint16(results))
			}
		}
	case "3xDint", "4xDint":
		if err := checkResultLength(results, 4); err != nil {
			readErr = err
		} else {
			adjustedResults, err := reverseBytes4(results, order)
			if err != nil {
				readErr = err
			} else {
				result = fmt.Sprintf("%d", int32(binary.BigEndian.Uint32(adjustedResults)))
			}
		}
	case "3xfloat", "4xfloat":
		if err := checkResultLength(results, 4); err != nil {
			readErr = err
		} else {
			adjustedResults, err := reverseBytes4(results, order)
			if err != nil {
				readErr = err
			} else {
				value := math.Float32frombits(binary.BigEndian.Uint32(adjustedResults))
				decValue := decimal.NewFromFloat(float64(value)).Round(3)
				result = decValue.String()
			}
		}
	case "3xbit", "4xbit":
		if err := checkResultLength(results, 2); err != nil {
			readErr = err
		} else {
			err := reverseBytes2(results, order)
			if err != nil {
				readErr = err
			} else {
				resultUint16 := binary.BigEndian.Uint16(results)
				data := make([]byte, 2)
				binary.BigEndian.PutUint16(data, resultUint16)
				result = int16(binary.LittleEndian.Uint16(data))
				var a uint16
				if v, ok := result.(int16); ok {
					a = uint16(v)
				} else {
					handleError(ctx, msg, fmt.Errorf("类型断言失败"), nil)
					return
				}
				bit, err := extractBit(a, extract)
				if err != nil {
					handleError(ctx, msg, fmt.Errorf("提取位时出错: %v", err), nil)
					return
				}
				result = fmt.Sprintf("%d", bit)
			}
		}
	default:
		handleError(ctx, msg, fmt.Errorf("不支持的类型"), nil)
		return
	}

	if readErr != nil {
		handleError(ctx, msg, fmt.Errorf("读取寄存器失败: %s", readErr.Error()), result)
		return
	}

	resultsMap := map[string]interface{}{"value": result}
	msgDataBytes, err := json.Marshal(resultsMap)
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("转换 JSON 时出错: %v", err), result)
		return
	}

	msg.SetData(string(msgDataBytes))
	ctx.TellSuccess(msg)
}

// ModbusTCP_Write 函数用于写入 Modbus TCP 协议的数据
func ModbusTCP_Write(ctx types.RuleContext, msg types.RuleMsg) {
	ip, port, id, typess, start, order, valueStr, extract, err := extractCommonFields(msg.GetData())
	if err != nil {
		handleError(ctx, msg, err, nil)
		return
	}

	handler := modbus.NewTCPClientHandler(ip + ":" + port)
	handler.Timeout = 10 * time.Second
	handler.SlaveId = id

	err = handler.Connect()
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("Connect failed: %s", err.Error()), nil)
		return
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	var valueBytes []byte

	switch typess {
	case "0xbit":
		if valueStr == "1" {
			valueBytes = []byte{0, 1}
		} else {
			valueBytes = []byte{0, 0}
		}
	case "4xint":
		valueInt, err := strconv.ParseInt(valueStr, 10, 16)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("解析 int16 值失败: %v", err), nil)
			return
		}
		valueBytes = make([]byte, 2)
		binary.BigEndian.PutUint16(valueBytes, uint16(valueInt))
	case "4xuint":
		valueUint, err := strconv.ParseUint(valueStr, 10, 16)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("解析 uint16 值失败: %v", err), nil)
			return
		}
		valueBytes = make([]byte, 2)
		binary.BigEndian.PutUint16(valueBytes, uint16(valueUint))
	case "4xDint":
		valueInt, err := strconv.ParseInt(valueStr, 10, 32)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("解析 int32 值失败: %v", err), nil)
			return
		}
		valueBytes = make([]byte, 4)
		binary.BigEndian.PutUint32(valueBytes, uint32(valueInt))
	case "4xfloat":
		valueFloat, err := strconv.ParseFloat(valueStr, 32)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("解析 float32 值失败: %v", err), nil)
			return
		}
		valueBytes = make([]byte, 4)
		binary.BigEndian.PutUint32(valueBytes, math.Float32bits(float32(valueFloat)))
	case "4xbit":
		results, err := client.ReadHoldingRegisters(start, 1)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("读取寄存器失败: %s", err.Error()), nil)
			return
		}
		if err := checkResultLength(results, 2); err != nil {
			handleError(ctx, msg, fmt.Errorf("读取寄存器失败: %s", err.Error()), nil)
			return
		}
		err = reverseBytes2(results, order)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("反转字节序失败: %s", err.Error()), nil)
			return
		}
		currentValue := binary.BigEndian.Uint16(results)
		bitValue, err := strconv.ParseUint(valueStr, 10, 1)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("解析位值失败: %v", err), nil)
			return
		}
		if bitValue == 1 {
			currentValue |= uint16(1 << extract)
		} else {
			currentValue &^= uint16(1 << extract)
		}
		valueBytes = make([]byte, 2)
		binary.LittleEndian.PutUint16(valueBytes, currentValue)
	default:
		handleError(ctx, msg, fmt.Errorf("不支持的类型"), nil)
		return
	}

	switch typess {
	case "4xbit":
		err = reverseBytes2(valueBytes, order)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("反转字节序失败: %s", err.Error()), nil)
			return
		}
		_, err = client.WriteMultipleRegisters(start, 1, valueBytes)
	case "4xint", "4xuint":
		err = reverseBytes2(valueBytes, order)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("反转字节序失败: %s", err.Error()), nil)
			return
		}
		valueBytes[0], valueBytes[1] = valueBytes[1], valueBytes[0]
		_, err = client.WriteMultipleRegisters(start, 1, valueBytes)
	case "4xDint", "4xfloat":
		valueBytes, err = reverseBytes4(valueBytes, order)
		if err != nil {
			handleError(ctx, msg, fmt.Errorf("反转字节序失败: %s", err.Error()), nil)
			return
		}
		_, err = client.WriteMultipleRegisters(start, 2, valueBytes)
	}

	if err != nil {
		handleError(ctx, msg, fmt.Errorf("写入寄存器失败: %s", err.Error()), nil)
		return
	}

	var valueBytesStr interface{}

	switch typess {
	case "0xbit":
		valueBytesStr = fmt.Sprintf("%d", binary.BigEndian.Uint16(valueBytes))
	case "4xbit":
		valueBytesStr = fmt.Sprintf("%d", binary.BigEndian.Uint16(valueBytes))
	case "4xint":
		valueBytesStr = fmt.Sprintf("%d", int16(binary.BigEndian.Uint16(valueBytes)))
	case "4xuint":
		valueBytesStr = fmt.Sprintf("%d", binary.BigEndian.Uint16(valueBytes))
	case "4xDint":
		valueBytesStr = fmt.Sprintf("%d", int32(binary.BigEndian.Uint32(valueBytes)))
	case "4xfloat":
		valueBytesStr = fmt.Sprintf("%.4f", math.Float32frombits(binary.BigEndian.Uint32(valueBytes)))
	default:
		handleError(ctx, msg, fmt.Errorf("不支持的类型"), nil)
		return
	}

	resultsMap := map[string]interface{}{"value": valueBytesStr}
	msgDataBytes, err := json.Marshal(resultsMap)
	if err != nil {
		handleError(ctx, msg, fmt.Errorf("转换 JSON 时出错: %v", err), nil)
		return
	}

	msg.SetData(string(msgDataBytes))
	fmt.Println("msg.Data:", msg.GetData())
	ctx.TellSuccess(msg)
}
