package concurrentmap

import (
	"errors"
	"fmt"
	"sync"
)

// 核心数据结构，map，key是车次：日期：类型，value是二维位图，x是长度，y是座位位图
const SHARD_COUNT = 32

// var myMap = New()

// func GetMap() ConcurrentMap {
// 	return myMap
// }

// // Map 分片
type ConcurrentMap []*ConcurrentMapShared

// 每一个Map 是一个加锁的并发安全Map
type ConcurrentMapShared struct {
	items        map[string]tickets
	sync.RWMutex // 各个分片Map各自的锁
}
type tickets struct {
	segments [][]uint8
}

//封装了读后写的逻辑
func (s tickets) Update(count, startNo, EndNo int32) ([]int64, error) {
	resBytes := make([]uint8, len(s.segments[0]))
	copy(resBytes, s.segments[0])
	fmt.Println("s.segments[0]", resBytes)
	//对每个区间作与运算
	for i := 1; i < len(s.segments); i++ {
		//计算Seats 经过与运算后的位图
		for j := 0; j < len(resBytes); j++ {
			resBytes[j] = resBytes[j] & s.segments[i][j]
		}
		fmt.Println("与运算后bytes", resBytes)
	}
	numbers, ok := findCount(resBytes, count)
	if ok != true {
		return nil, errors.New("余票不足")
	} else {
		fmt.Println("numbers", numbers)
		// update()
		for i := 0; i < len(numbers); i++ {
			setZero(s.segments, numbers[i])
		}
		fmt.Println("s.segments", s.segments)
		return numbers, nil
	}
}

//findCount
func findCount(bytes []uint8, count int32) ([]int64, bool) {
	var res []int64
	num := int64(1)
	k := int32(0)
	for i := 0; i < len(bytes); i++ {
		x := uint8(128)
		temp := bytes[i]
		for x != 0 {
			if temp&x != 0 {
				res = append(res, num)
				k++
				if k == count {
					return res, true
				}
			}
			num++
			x = x >> 1
		}
	}
	return nil, false
}
func setOne(bytes [][]uint8, validSeatNo int64) {
	index := (validSeatNo - 1) >> 3
	pos := (validSeatNo - 1) % 8
	//对每个区间更改
	for i := 0; i < len(bytes); i++ {
		bytes[i][index] = bytes[i][index] | (128 >> pos)
	}
}
func setZero(bytes [][]uint8, validSeatNo int64) {

	index := (validSeatNo - 1) >> 3
	pos := (validSeatNo - 1) % 8
	//对每个区间更改
	for i := 0; i < len(bytes); i++ {
		bytes[i][index] = bytes[i][index] & ^(128 >> pos)
	}
}
func NewConcurrentMap() ConcurrentMap {
	// SHARD_COUNT 默认32个分片
	m := make(ConcurrentMap, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &ConcurrentMapShared{
			items: make(map[string]tickets),
		}
	}
	return m
}

func (m ConcurrentMap) GetShard(key string) *ConcurrentMapShared {
	return m[uint(fnv32(key))%uint(SHARD_COUNT)]
}

// FNV hash
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// func (m ConcurrentMap) Set(key string, value tickets) {
// 	shard := m.GetShard(key) // 段定位找到分片
// 	shard.Lock()             // 分片上锁
// 	shard.items[key] = value // 分片操作
// 	shard.Unlock()           // 分片解锁
// }
func (m ConcurrentMap) Get(key string) (tickets, bool) {
	shard := m.GetShard(key)
	// shard.RLock()
	val, ok := shard.items[key]
	// shard.RUnlock()
	return val, ok
}
func (m ConcurrentMap) Init(key string, tripLen, seats int) {
	shard := m.GetShard(key) // 段定位找到分片
	shard.Lock()             // 分片上锁
	t := tickets{}
	t.segments = make([][]uint8, tripLen)
	a := []uint8{255, 255}
	for i := 0; i < tripLen; i++ {
		t.segments[i] = a
	}
	// fmt.Println(t)
	shard.items[key] = t // 分片操作
	shard.Unlock()       // 分片解锁
}
func (m ConcurrentMap) Update(key string, count, startNo, EndNo int32) ([]int64, error) {
	shard := m.GetShard(key) // 段定位找到分片
	shard.Lock()             // 分片上锁
	defer shard.Unlock()     // 分片解锁
	val, ok := shard.items[key]
	fmt.Print(val)
	if ok == false {
		return nil, errors.New("不存在该key，无法更新")
	} else {
		return val.Update(count, startNo, EndNo)
	}
}
