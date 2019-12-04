package util

import (
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
