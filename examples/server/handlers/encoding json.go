package handlers

import (
	"encoding/json"
	"fmt"
	"log"
)

// 定义一个函数来把订阅的json数据转换为map,用于提取数据指定字段的值，是给其它程序调用的函数
func TransformData(subscriptionName string, msgIn []byte) (msgOut map[string]interface{}, err error) {
	msgOut = make(map[string]interface{}) // 初始化 msgOut

	// 尝试解析消息内容为 JSON
	var result interface{}
	if err := json.Unmarshal(msgIn, &result); err != nil {
		log.Printf("解析 JSON 时出错: %v", err)
		return nil, fmt.Errorf("解析 JSON 时出错: %w", err)
	}

	// 定义一个递归函数来查找所有键值对
	var findName func(data interface{}, path []string, pathStr string)
	findName = func(data interface{}, path []string, pathStr string) {
		// 根据数据类型选择不同的处理方式
		switch v := data.(type) {
		case map[string]interface{}:
			// 处理单个对象
			for k, val := range v {
				newPath := append(path, k)
				newPathStr := pathStr + "." + k
				// 判断是否是目标键
				if k == subscriptionName {
					// 打印键值对应的值
					msgOut[newPathStr[1:]] = val // 去掉开头的点
					//fmt.Printf("找到键 '%s' 的值: %v\n", newPathStr[1:], val)//打印结果对应的值
				}
				findName(val, newPath, newPathStr)
			}
		case []interface{}:
			// 处理数组
			for i, item := range v {
				newPath := append(path, fmt.Sprintf("[%d]", i))
				newPathStr := pathStr + fmt.Sprintf(".[%d]", i)
				findName(item, newPath, newPathStr)
			}
		}
	}

	//fmt.Printf("开始查找 JSON 中的键值对\n")//打印开始查找的提示
	findName(result, []string{}, "") // 开始查找 JSON 中的键值对
	//fmt.Printf("查找完成\n")//打印查找完成的提示

	return msgOut, nil //返回结果和错误信息
}
