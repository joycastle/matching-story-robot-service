package config

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ToRangeIndexLCR8(data map[int]int) map[string]int {
	maxIndex := 0
	last := -1
	ret := make(map[string]int)
	for k, v := range data {
		min := last + 1
		max := k
		last = k

		key := fmt.Sprintf("%d-%d", min, max)
		ret[key] = v

		if maxIndex < v {
			maxIndex = v
		}
	}

	key := fmt.Sprintf("%d-%d", last+1, math.MaxInt)
	ret[key] = maxIndex + 1

	return ret
}

func RangeIndexLORC(data map[string]int, level int) int {
	defaultIndex := 0
	i := 0
	for k, index := range data {
		if i == 0 {
			defaultIndex = index
		}
		i++
		arr := strings.Split(k, "-")
		min, _ := strconv.Atoi(arr[0])
		max, _ := strconv.Atoi(arr[1])
		if level > min && level <= max {
			return index
		}
	}
	return defaultIndex
}

func Compare2Int(a1, a2 int) (int, int) {
	if a1 > a2 {
		return a2, a1
	} else {
		return a1, a2
	}
}

func ParseStringType(v string) [][]int {
	v = strings.TrimLeft(v, "{")
	v = strings.TrimRight(v, "}")
	arr := strings.Split(v, ",")
	ret := [][]int{}
	for _, v := range arr {
		arrr := strings.Split(v, "|")
		tmp := []int{}
		for _, vv := range arrr {
			vvInt, _ := strconv.Atoi(vv)
			tmp = append(tmp, vvInt)
		}
		ret = append(ret, tmp)
	}
	return ret
}
