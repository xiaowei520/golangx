package main

import (
	"strings"
	"testing"
)

//BenchmarkCalculateIp-8   	 5000000	       289 ns/op
//1ms=100W ns
//6000ms ---    600KW ns
//1,000,000 纳秒 = 1毫秒
//IP转换 5s 多
func BenchmarkCalculateIp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateIp("1.1.1.1")
	}

}

//goarch: amd64
//pkg: gitlab.meiyou.com/biz-modules/timing-push/conf
//BenchmarkReadLine-8   	200000000	         8.54 ns/op
func BenchmarkReadLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Index("122.226.129.25	-	-	30/Sep/2019:08:57:10 +0800	-	77500	11750	4	32040	POST	api.immomo.com	/v2/setting/abtest/index?fr=746208866	HTTP/2.0	691	0.040	200540	549	-	MomoChat/8.20.1 Android/4948 (SEA-AL10; Android 9; Gapps 0; zh_CN; 17; HUAWEI)	32070492501	22	10.223.21.3:9000	200	346	0.040", "-")
	}
}

func BenchmarkReadLine2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Read([]byte("122.226.129.25	-	-	30/Sep/2019:08:57:10 +0800	-	77500	11750	4	32040	POST	api.immomo.com	/v2/setting/abtest/index?fr=746208866	HTTP/2.0	691	0.040	200540	549	-	MomoChat/8.20.1 Android/4948 (SEA-AL10; Android 9; Gapps 0; zh_CN; 17; HUAWEI)	32070492501	22	10.223.21.3:9000	200	346	0.040"))
	}

}
func Read(line []byte) int {
	var result int
	for i := 7; i <= 15; i++ {
		if line[i] == ' ' || line[i] == '-' {
			result = i
			break
		}
	}
	return result
}
