package rate

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

//使用cas 锁来实现、减少开销
//业务限流器
const (
	Normal = iota
	AbNormal
)

//标示, 可以是appId  uid 唯一性的代表值
type Limit float64

type Limiter struct {
	limit     Limit
	EventName string
	//cas
	done uint32

	//时间
	retryTime int
	//重试次数
	retryLen int

	//最大缓冲时间
	maxTime int64
	//最小缓冲时间
	minTime int64

	mu sync.Mutex
	// last is the last time the limiter's tokens field was updated
	last time.Time
	// lastEvent is the latest time of a rate-limited event (past or future)
	lastEvent time.Time
}

func NewLimiter(r Limit, n string) *Limiter {

	return &Limiter{
		limit:     r,
		EventName: n,
	}
}

func DefaultLimiter(r Limit, n string) *Limiter {
	return &Limiter{
		limit:     r,
		EventName: n,
		maxTime:   1000,
		minTime:   100,
	}
}

func (lim *Limiter) Allow() bool {
	//no limit rate
	if atomic.LoadUint32(&lim.done) == Normal {
		return true
	}
	//rand sleep, avoid conflict
	randTime := rand.Int63n(lim.maxTime-lim.minTime) + lim.minTime
	time.Sleep(time.Duration(randTime) * time.Millisecond)
	return false
}

//where use stop
func (lim *Limiter) Stop() bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	if lim.done == Normal {
		defer atomic.StoreUint32(&lim.done, AbNormal)
	}

	return true
}

func (lim *Limiter) Recover() bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	if lim.done == AbNormal {
		defer atomic.StoreUint32(&lim.done, Normal)
	}
	return true
}
