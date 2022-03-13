package util

import (
	"crypto/rand"
	"github.com/shopspring/decimal"
	"math/big"
	mathRand "math/rand"
)

func GetRandom(min, max int64) int64 {
	result, _ := rand.Int(rand.Reader, big.NewInt(max))
	resultInt := result.Int64()
	return resultInt%(max-min+1) + min
}

func RandomRangeSlice(nums []int64, factor float64, rangeNum float64) []int64 {
	rangedNums := make([]int64, len(nums))
	for i := 0; i < len(nums); i++ {
		rangeNumTemp := rangeNum
		if nums[i] < 10 {
			rangeNumTemp = rangeNum * 3
		} else if nums[i] < 100 {
			rangeNumTemp = rangeNum * 2
		}
		numDec := decimal.NewFromInt(nums[i]).Mul(decimal.NewFromFloat(factor))
		min := numDec.Mul(decimal.NewFromFloat(1 - rangeNumTemp)).IntPart()
		max := numDec.Mul(decimal.NewFromFloat(1 + rangeNumTemp)).IntPart()
		random := GetRandom(min, max)
		rangedNums[i] = random
	}
	return rangedNums
}

func ShuffledSlice(nums []int64) []int64 {
	numsTemp := make([]int64, len(nums))
	copy(numsTemp, nums)
	shuffled := make([]int64, len(numsTemp))
	for i := range shuffled {
		j := mathRand.Intn(len(numsTemp))
		shuffled[i] = numsTemp[j]
		numsTemp = append(numsTemp[:j], numsTemp[j+1:]...)
	}
	return shuffled
}

func ExchangeChosenIndex(nums []int64, chosenNum int64, chosenIndex int) []int64 {
	numsTemp := make([]int64, len(nums))
	copy(numsTemp, nums)
	for i := range numsTemp {
		if numsTemp[i] == chosenNum && i != chosenIndex {
			temp := numsTemp[chosenIndex]
			numsTemp[chosenIndex] = chosenNum
			numsTemp[i] = temp
			break
		}
	}
	return numsTemp
}
