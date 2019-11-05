package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//文章地址
//https://segmentfault.com/a/1190000020916113

const N = 256

//构建N个堆
var GlobalIp map[int64]*IpQueue

//然后N个堆 获取TOP10
var GlobalNum map[int64]int64 //次数

//思路1: ?- 是否可以组合我们的超大数字- 组合方式 出现次数+ 十进制数字、 堆排序、直接就能得到结果集-好处是避免了反射、

//思路2: 2.1 直接将IP变成十进制 hash算次数。2.2 mod N 进行堆排序 2.3 进行N个堆TOP10 排序聚合 | 2.4 输出聚合后的堆TOP10

func ReadLine(filePth string, hookfn func([]byte)) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		hookfn(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

}

//
func initHeap() {
	GlobalNum = make(map[int64]int64)
	GlobalIp = make(map[int64]*IpQueue)
	for i := 0; i <= N; i++ {
		q := make(IpQueue, 1)
		q[0] = &Item{ip: "0.0.0.0", num: -1}
		heap.Init(&q)
		GlobalIp[int64(i)] = &q //堆给到全局Global
	}
}

//2.1 直接将IP变成十进制 hash算次数
func processLine(line []byte) {

	var result int
	for i := 7; i <= 15; i++ {
		if line[i] == '\t' || line[i] == '-' {
			result = i
			break
		}
	}
	str := string(line[0:result])

	ipv4 := CalculateIp(string(str))

	GlobalNum[int64(ipv4)]++
}

//2.2 mod N 进行堆排序
func handleHash() {
	//堆耗时开始
	timestamp := time.Now().UnixNano() / 1000000
	for k, v := range GlobalNum {
		heap.Push(GlobalIp[k%N], &Item{ip: RevIp(k), num: int64(v)})
	}
	edgiest := time.Now().UnixNano() / 1000000
	fmt.Println("堆耗时总时间ms:", edgiest-timestamp)
}

//2.3 进行N个堆TOP10 排序聚合
func polyHeap() {
	//聚合N 个 小堆的top10
	for i := 0; i < N; i++ {
		iterator := 10
		if iterator > GlobalIp[int64(i)].Len() {
			iterator = GlobalIp[int64(i)].Len()
		}
		for j := 0; j < iterator; j++ {
			//写入到堆栈N
			item := heap.Pop(GlobalIp[int64(i)]).(*Item)
			heap.Push(GlobalIp[N], item)
		}
	}
}

//2.4 输出聚合后的堆TOP10
func printResult() {
	result := 0
	for result < 10 {
		item := heap.Pop(GlobalIp[N]).(*Item)
		fmt.Printf("出现的次数:%d|IP:%s \n", item.num, item.ip)
		result++
	}
}

//string 转IP
func CalculateIp(str string) int64 {
	x := strings.Split(str, ".")
	b0, _ := strconv.ParseInt(x[0], 10, 0)
	b1, _ := strconv.ParseInt(x[1], 10, 0)
	b2, _ := strconv.ParseInt(x[2], 10, 0)
	b3, _ := strconv.ParseInt(x[3], 10, 0)

	number0 := b0 * 16777216 //256*256*256
	number1 := b1 * 65536    //256*256
	number2 := b2 * 256      //256
	number3 := b3 * 1        //1
	sum := number0 + number1 + number2 + number3
	return sum
}

//ip 转string
func RevIp(ip int64) string {

	ip0 := ip / 16777216 //高一位
	ip1 := (ip - ip0*16777216) / 65536
	ip2 := (ip - ip0*16777216 - ip1*65536) / 256
	ip3 := ip - ip0*16777216 - ip1*65536 - ip2*256
	return fmt.Sprintf("%d.%d.%d.%d", ip0, ip1, ip2, ip3)
}

type Item struct {
	ip  string
	num int64
}

type IpQueue []*Item

func (pq IpQueue) Len() int { return len(pq) }

func (pq IpQueue) Less(i, j int) bool {
	return pq[i].num > pq[j].num
}
func (pq IpQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *IpQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *IpQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func main() {
	runtime.GOMAXPROCS(2)
	timestamp := time.Now().UnixNano() / 1000000

	//初始化
	initHeap()

	//串行 读取文件 写入到hash map
	_ = ReadLine("/Users/admin/Downloads/api.immomo.com-access_10-01.log", processLine)

	//多个小堆
	handleHash()

	//聚合堆
	polyHeap()

	//打印结果

	printResult()

	fmt.Println(time.Now().UnixNano()/1000000 - timestamp)
}
