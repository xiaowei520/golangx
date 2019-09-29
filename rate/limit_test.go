package rate

import (
	"fmt"
	"testing"
)

func TestNewLimiter(t *testing.T) {
	test := DefaultLimiter(1, "test")

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

//BenchmarkLimiter_Allow-8   	     100	 546062919 ns/op  当不被恢复的时候、大约500ms sleep
//BenchmarkLimiter_Allow-8   	200000000	         6.69 ns/op  当被恢复的时候 6.69ns 每次
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
