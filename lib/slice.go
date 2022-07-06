package lib

func ArraySliceInt64(ids []int64, size int) [][]int64 {
	var ret [][]int64

	if len(ids) == 0 {
		return ret
	}

	count := 0
	tmp := []int64{}
	for _, id := range ids {
		tmp = append(tmp, id)
		count++
		if count == size {
			ret = append(ret, tmp)
			count = 0
			tmp = []int64{}
		}
	}

	if count > 0 {
		ret = append(ret, tmp)
	}

	return ret
}
