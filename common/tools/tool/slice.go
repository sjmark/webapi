package tool

//SliceDiff 数组差集
func SliceDiff(val, vel []int64) []int64 {
	diffSlice := make([]int64, 0, len(val))

	for j := 0; j < len(val); j++ {
		for i := 0; i < len(vel); i++ {
			if val[j] == vel[i] {
				break
			} else {
				if len(vel) == (i + 1) { //如果不同 查看a的元素个数及当前比较元素的位置 将不同的元素添加到返回slice中
					diffSlice = append(diffSlice, val[j])
				}
			}
		}
	}
	return diffSlice
}

//SliceFlip 数组翻转s
func SliceFlip(val []int64) {

	for i := 0; i < len(val)/2; i++ {
		val[i], val[len(val)-1-i] = val[len(val)-1-i], val[i]
	}
}
