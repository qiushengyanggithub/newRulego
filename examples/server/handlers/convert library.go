package handlers

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ReverseBytes2 根据指定顺序调整2字节的顺序，给MODBUS读取数据时使用
func reverseBytes2(data []byte, order string) error {
	if len(data) != 2 {
		return fmt.Errorf("数据长度必须为2字节")
	}

	switch order {
	case "BA":
		// 保持不变，不交换字节
		return nil
	case "AB":
		// 直接交换字节
		data[0], data[1] = data[1], data[0]
		return nil
	default:
		return fmt.Errorf("不支持的顺序: %s", order)
	}
}

// ReverseBytes4 根据指定顺序调整4字节的顺序
func reverseBytes4(data []byte, order string) ([]byte, error) {
	if len(data) != 4 {
		return data, fmt.Errorf("数据长度必须为4字节")
	}

	switch order {
	case "ABCD", "AB", "BA":
		return data, nil
	case "CDAB":
		return []byte{data[2], data[3], data[0], data[1]}, nil
	case "BADC":
		return []byte{data[1], data[0], data[3], data[2]}, nil
	case "DCBA":
		return []byte{data[3], data[2], data[1], data[0]}, nil
	default:
		return data, fmt.Errorf("不支持的顺序: %s", order)
	}
}

func extractBit(value uint16, bitValue uint16) (uint16, error) { //从uint16值中提取指定位的值，从右边往左边 0-15确定位位置
	position := bitValue

	if position < 0 || position > 15 {
		return 0, fmt.Errorf("bit位置必须在0到15之间")
	}

	// 使用位操作提取指定位置的bit
	bit := (value >> position) & 1
	return bit, nil
}

func BytesToInt16(b []byte) int16 { //将字节切片转换为int16
	return int16(binary.BigEndian.Uint16(b))
}
func BytesToUint16(b []byte) uint16 { //将字节切片转换为uint16
	return binary.BigEndian.Uint16(b)
}
func BytesToInt32(b []byte) int32 { //将字节切片转换为int32
	return int32(binary.BigEndian.Uint32(b))
}
func BytesToUint32(b []byte) uint32 { //将字节切片转换为uint32
	return binary.BigEndian.Uint32(b)
}
func BytesToFloat32(b []byte) float32 { //将字节切片转换为float32
	return math.Float32frombits(binary.BigEndian.Uint32(b))
}

func BytesToFloat64(b []byte) float64 { //将字节切片转换为float64
	return math.Float64frombits(binary.BigEndian.Uint64(b))
}

// String  将传入的类型转为字符串
func String(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	case int:
		return strconv.FormatInt(int64(val.(int)), 10)
	case int8:
		return strconv.FormatInt(int64(val.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(val.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(val.(int32)), 10)
	case int64:
		return strconv.FormatInt(int64(val.(int64)), 10)
	case uint:
		return strconv.FormatUint(uint64(val.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(val.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(val.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(val.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(uint64(val.(uint64)), 10)
	case bool:
		return strconv.FormatBool(val.(bool))
	case float32:
		return strconv.FormatFloat(val.(float64), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', -1, 64)
	default:
		return ""
	}
}

// FormatFloat 浮点数转String 保留指定位数的小数
func FormatFloat(v interface{}, decimal int) string {
	switch v.(type) {
	case float32:
		return strconv.FormatFloat(float64(v.(float32)), 'f', decimal, 32)
	case float64:
		return strconv.FormatFloat(float64(v.(float64)), 'f', decimal, 64)
	default:
		return ""
	}
}

// FormatInt 整数转字符串
func FormatInt(val interface{}) string {
	switch val.(type) {
	case int:
		return strconv.FormatInt(int64(val.(int)), 10)
	case int8:
		return strconv.FormatInt(int64(val.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(val.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(val.(int32)), 10)
	case int64:
		return strconv.FormatInt(int64(val.(int64)), 10)
	case uint:
		return strconv.FormatUint(uint64(val.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(val.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(val.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(val.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(uint64(val.(uint64)), 10)
	default:
		return ""
	}
}

// Md5 对指定的字符串进行MD5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	md5Byte := h.Sum(nil)
	md5Str := hex.EncodeToString(md5Byte)
	return md5Str
}

// Md5Upper 返回大写的MD5十六进制加密后字符串
func Md5Upper(str string) string {
	return strings.ToUpper(Md5(str))
}

// Md5Lower 返回小写的MD5十六进制加密后字符串
func Md5Lower(str string) string {
	return strings.ToLower(Md5(str))
}

// SHA1 对给定的字符串进行SHA1加密
func SHA1(str string) string {
	s := sha1.New()
	s.Write([]byte(str))
	md5Byte := s.Sum(nil)
	sha1Str := hex.EncodeToString(md5Byte)
	return sha1Str
}

// SHA512 对给定的字符串进行SHA512加密
func SHA512(str string) string {
	s := sha512.New()
	s.Write([]byte(str))
	md5Byte := s.Sum(nil)
	sha1Str := hex.EncodeToString(md5Byte)
	return sha1Str
}

// SHA256 对给定的字符串进行SHA256加密
func SHA256(str string) string {
	s := sha256.New()
	s.Write([]byte(str))
	md5Byte := s.Sum(nil)
	sha1Str := hex.EncodeToString(md5Byte)
	return sha1Str
}
