package visitor

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Visitor struct {
	Limiter            *rate.Limiter
	LastSeen           time.Time
	ReputationScore    int
	IsBlocked          bool
	IsPermanentBlocked bool
	BlockUntil         time.Time
	LastRateLimitTime  time.Time
	NotFoundCount      int
}

var visitors = make(map[string]*Visitor)
var mu sync.Mutex

func GetVisitor(ip string) *Visitor {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]

	if !exists {
		v = &Visitor{
			Limiter: rate.NewLimiter(1, 5),
		}
		visitors[ip] = v
	}

	v.LastSeen = time.Now()

	return v
}

func CleanupVisitors() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			<-ticker.C
			mu.Lock()
			fmt.Println("Cleaning up expired visitors list")
			for ip, visitor := range visitors {
				// if time.Now().After(visitor.lastSeen.Add(3 * time.Minute)) {
				// 	fmt.Println("Deleting expired visitor:", ip)
				// 	delete(visitors, ip)
				// }

				if !visitor.IsPermanentBlocked {
					if visitor.IsBlocked {
						if time.Now().After(visitor.BlockUntil) {
							mu.Lock()
							visitor.IsBlocked = false
							if visitor.ReputationScore > 0 {
								visitor.ReputationScore--
							}
							mu.Unlock()
						}
					} else {
						if time.Now().After(visitor.LastRateLimitTime.Add(10 * time.Minute)) {
							fmt.Printf("Resetting reputation score for IP %s\n", ip)
							mu.Lock()
							if visitor.ReputationScore > 0 {
								visitor.ReputationScore--
							}
							mu.Unlock()
						}
					}
				}
			}
			mu.Unlock()
		}
	}()
}
