package middlewares

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/httpcommon"
	"golang.org/x/time/rate"
)

type userClient struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type UserRateLimiter struct {
	clients map[string]*userClient
	mu      sync.Mutex
	r       rate.Limit
	burst   int

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewUserRateLimiter(r rate.Limit, burst int) *UserRateLimiter {
	ctx, cancel := context.WithCancel(context.Background())

	rl := &UserRateLimiter{
		clients: make(map[string]*userClient),
		r:       r,
		burst:   burst,
		ctx:     ctx,
		cancel:  cancel,
	}

	rl.wg.Add(1)
	go rl.cleanupLoop()
	return rl
}

func (rl *UserRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		local, err := httpcommon.GetLocalSession(r)

		if err != nil {
			httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(unauthorized))
			return
		}

		limiter := rl.getLimiter(local.UserID)
		if !limiter.Allow() {
			httpcommon.WriteJSON(w, http.StatusTooManyRequests, httpcommon.Message(tooManyRequests))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *UserRateLimiter) getLimiter(userID string) *rate.Limiter {

	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[userID]
	if !exists {
		lim := rate.NewLimiter(rl.r, rl.burst)
		rl.clients[userID] = &userClient{limiter: lim, lastSeen: time.Now()}
		return lim
	}

	c.lastSeen = time.Now()
	return c.limiter
}

func (rl *UserRateLimiter) cleanupLoop() {
	defer rl.wg.Done()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-rl.ctx.Done():
			return
		case <-ticker.C:
			rl.mu.Lock()
			for userID, c := range rl.clients {
				if time.Since(c.lastSeen) > 10*time.Minute {
					delete(rl.clients, userID)
				}
			}
			rl.mu.Unlock()
		}
	}
}

func (rl *UserRateLimiter) Stop() {
	rl.cancel()
	rl.wg.Wait()
}
