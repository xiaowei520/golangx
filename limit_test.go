package core

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
