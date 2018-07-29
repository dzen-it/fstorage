package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Throttler struct {
	visitors map[string]*visitor
	rps      float64
	cs       int32
	mx       sync.Mutex
}

//NewThrottler sets:
//request per second
//num of concurent sessions
func NewThrottler(rps float64, cs int) *Throttler {
	t := &Throttler{
		visitors: make(map[string]*visitor),
		rps:      rps,
		cs:       int32(cs),
		mx:       sync.Mutex{},
	}

	go cleanerWorker(t)
	return t
}

type visitor struct {
	rpsLimiter *rate.Limiter
	csLimiter  int32
	lastSeen   time.Time
}

func (t *Throttler) getVisitor(ip string) *visitor {
	v, ok := t.visitors[ip]
	if !ok {
		v = &visitor{
			rpsLimiter: rate.NewLimiter(rate.Limit(t.rps), 1),
			csLimiter:  0,
		}
		t.visitors[ip] = v
	}
	v.lastSeen = time.Now()
	return v
}

func (t *Throttler) IsAllow(ip string) bool {
	t.mx.Lock()
	defer t.mx.Unlock()

	v := t.getVisitor(ip)

	if !v.rpsLimiter.Allow() {
		return false
	}
	if v.csLimiter > t.cs {
		return false
	}

	return true
}

func (t *Throttler) AddSession(ip string) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.getVisitor(ip).csLimiter++
}

func (t *Throttler) CloseSession(ip string) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.getVisitor(ip).csLimiter--
}

func cleanerWorker(t *Throttler) {
	for {
		time.Sleep(time.Minute)
		t.mx.Lock()
		for ip, v := range t.visitors {
			if time.Now().Sub(v.lastSeen) > 3*time.Minute {
				delete(t.visitors, ip)
			}
		}
		t.mx.Unlock()
	}
}

func MiddlewareThrottling(t *Throttler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rip, _ := net.ResolveTCPAddr("tcp", r.RemoteAddr)
			ip := rip.IP.String()
			ok := t.IsAllow(ip)
			if !ok {
				http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
				return
			}

			t.AddSession(ip)
			next.ServeHTTP(w, r)
			t.CloseSession(ip)
		})
	}
}
