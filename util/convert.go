package util

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"
)

// 将字符串数组解析成数字数组
func AtoInt(arr []string) ([]int, error) {
	if len(arr) == 0 {
		return []int{}, nil
	}

	result := make([]int, len(arr))
	var err error
	for i, v := range arr {
		if result[i], err = strconv.Atoi(v); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// 将以逗号隔开的字符串数字分解
func AtoInt2(str string) ([]int, error) {
	str = strings.ReplaceAll(str, " ", "")
	arr := strings.Split(str, ",")
	return AtoInt(arr)
}

// 随机生成前缀为 prefix、总长度为 length 的 Token，长度不足时，追加随机字符
func RandToken(prefix string, length int) string {
	chars := "01234567890abcdefghijklmnopqrstuvwxyz"
	appendLength := length - len(prefix)

	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(prefix)

	for i := 0; i < appendLength; i++ {
		index := rand.Intn(len(chars))
		buffer.WriteString(chars[index : index+1])
	}
	return buffer.String()
}

// 生成随机端口
func RandPort() int {
	return 10000 + rand.Intn(65535-10000)
}

// 批量生成随机端口
func RandPorts(count int) []int {
	result := make([]int, count)
	result[0] = RandPort()

	for i := 1; i < count; i++ {
		result[i] = result[i-1] + 1
	}
	return result
}
