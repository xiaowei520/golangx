package rate

import (
	"fmt"
	"testing"
)

func TestNewLimiter(t *testing.T) {
	test := DefaultLimiter(1, "test")
	test1 := DefaultLimiter(1, "test")
	test2 := DefaultLimiter(1, "test")
	test3 := DefaultLimiter(1, "test")

	fmt.Println(&test1)
	fmt.Println(&test2)
	fmt.Println(&test3)
	//fmt.Println(*DefaultLimiter(1,"3"))
	//fmt.Println(*DefaultLimiter(1,"4"))
	go func() {
		test.Stop()
	}()

	for i := 0; i < 100; i++ {
		fmt.Println(test.Allow())
		if i == 50 {
			test.Recover()
		}
	}

}

func TestLimiter_WindowAllow(t *testing.T) {
	test := SWLimiter(1, "test", 100000, 1)
	//fmt.Println(test.WindowAllow())
	for i := 0; i < 30000000; i++ {
		//time.Sleep(1 * time.Second)
		if test.WindowAllow() {
			//fmt.Println("true")
		}
	}
	fmt.Println(test)

}
func BenchmarkLimiter_WindowAllow(b *testing.B) {
	//BenchmarkLimiter_WindowAllow-8   	 3000000	       525 ns/op
	//test := SWLimiter(1, "test", 100000, 1)

	//BenchmarkLimiter_WindowAllow-8   	 3000000	       508 ns/op
	test := SWLimiter(1, "test", 1, 1)

	for i := 0; i < b.N; i++ {
		test.WindowAllow()
	}
}

//BenchmarkLimiter_Allow-8   	     100	 546062919 ns/op  当不被恢复的时候、大约500ms sleep
//BenchmarkLimiter_Allow-8   	200000000	         6.69 ns/op  当被恢复的时候 6.69ns 每次
//BenchmarkLimiter_Allow-8   	1000000000	         2.54 ns/op  当只进行allow时 且增加pad.
//BenchmarkLimiter_Allow-8   	1000000000	         3.00 ns/op   当不增加pad 时
func BenchmarkLimiter_Allow(b *testing.B) {
	test := DefaultLimiter(1, "test")
	for i := 0; i < b.N; i++ {
		test.Allow()
		if i == 2 {
			test.Stop()
		}
		if i == 3 {
			test.Recover()
		}
	}
}
