package main

import (
	"fmt"
	"luzijia/1368-i/cache"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	var testArr = []int64{5, 8, 10, 1, 3}
	fmt.Println(theMostDiff(testArr))
	//我用数字来代表54张牌
	var cards = []int{}
	for i := 0; i < 54; i++ {
		cards = append(cards, i)
	}
	fmt.Println(shuffle(cards))
	var people = []int{1, 2, 3, 4, 5}
	fmt.Println(game(3, 2, people))
	testMyCache()
}

/*
设计一个带失效时间的缓存数据结构，key和value都是string，并实现增删改查接口。
*/
func testMyCache() {
	cache1 := cache.NewExpiredMap()
	for i := 1; i <= 10; i++ {
		cache1.Set(i, strconv.Itoa(i), int64(i))
	}
	cache1.Delete(8)
	cache1.Delete(9)
	cache1.Delete(10)
	for i := 1; i <= 10; i++ {
		cache1.Set(10+i, strconv.Itoa(8), 5)
	}
	cache1.Close()
	time.Sleep(time.Millisecond)
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 1)
		cache1.Set(i, strconv.Itoa(i), int64(i))
		fmt.Println(cache1.Get(i))
	}
}

/*
实现一个游戏算法：输入n和m，代表有n个选手围成一圈（选手编号为0至n-1），0号从1开始报数，
报m的选手游戏失败从圆圈中退出，下一个人接着从1开始报数，如此反复，求最后的胜利者编号。
例如，n=3，m=2，那么失败者编号依次是1、0，最后的胜利者是2号。
这里考虑m，n都是正常的数据范围，其中：
1 <= n <= 10^5
1 <= m <= 10^6
算法要求考虑时间效率。
*/

func game(n, m int, nums []int) []int {
	return []int{1, 1}
}

/*
1. 实现一个函数，输入为任意长度的int64数组，返回元素最大差值，
例如输入arr=[5,8,10,1,3]，返回9。
*/
func theMostDiff(arr []int64) int64 {
	var res int64
	res = 0
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			t := arr[j] - arr[i]
			if t < 0 {
				t = -t
			}
			if t > res {
				res = t
			}
		}
	}
	return res
}

/*
实现一个函数，对输入的扑克牌执行洗牌，保证其是均匀分布的，
也就是说列表中的每一张扑克牌出现在列表的每一个位置上的概率必须相同。
*/
func shuffle(cards []int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]int, len(cards))
	n := len(cards)
	for i := 0; i < n; i++ {
		randIndex := r.Intn(len(cards))
		ret[i] = cards[randIndex]
		cards = append(cards[:randIndex], cards[randIndex+1:]...)
	}
	return ret
}
