package limit

import (
	"golang.org/x/time/rate"
	"net/http"
	"encoding/json"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func RateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
    limiter := rate.NewLimiter(2, 4)
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            message := Message{
                Status: "Request Failed",
                Body:   "The API is at capacity, try again later.",
            }

            w.WriteHeader(http.StatusTooManyRequests)
            json.NewEncoder(w).Encode(&message)
            return
        } else {
            next(w, r)
        }
    })
}