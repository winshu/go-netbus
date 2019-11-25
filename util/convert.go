package util

import (
	"bytes"
	"math/rand"
	"strconv"
)

func Atoi(arr []string) ([]int, error) {
	if len(arr) == 0 {
		return []int{}, nil
	}

	result := make([]int, len(arr))
	var err error
	for i, v := range arr {
		result[i], err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

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
