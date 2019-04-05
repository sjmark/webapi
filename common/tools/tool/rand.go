package tool

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//范围内随机 [min,max) 会取到min,不会取到max
func RandMinMaxForFloat(min float64, max float64) float64 {
	return min + (rand.Float64() * (max - min))
}

//范围内随机,[min,max] ,会取到边界值
func RandMinMaxForInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

//GetRandomInt 获取指定长度的随机数
func GetRandomInt(length int) string {
	var chars = make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		index := r.Intn(len(verifyrandom))
		if index == len(verifyrandom)-1 {
			index = 1
		}
		chars[i] = verifyrandom[index]
	}
	return BytesToStr(chars)
}

//GetRandomStr 获取指定长度的随机字符串
func GetRandomStr(length int) string {
	var sidchars = make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		index := r.Intn(len(sidrandom))
		sidchars[i] = sidrandom[index]
	}
	return BytesToStr(sidchars)
}

//RandomNum 获取随机数
func RandomNum(total int) string {
	avatars := make([]string, total-1)
	outavatars := make([]string, total-1)
	for i := 2; i < (total + 1); i++ {
		avatars[i-2] = strconv.Itoa(i)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	end := total - 2
	for i := 0; i < (total - 1); i++ {
		num := r.Intn(end + 1)
		outavatars[i] = avatars[num]
		avatars[num] = avatars[end]
		end = end - 1
	}
	res := strings.Join(outavatars, ",")
	return res
}

// 不重复的随机数
func RandomNumber(start int, count int) []int {
	//范围检查
	if start < 0 || (start) < count {
		return nil
	}
	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn(start)
		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}

func RandGenerator(n int) int {
	return rand.Intn(n)
}

func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}