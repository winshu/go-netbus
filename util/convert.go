package util

import "strconv"

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
