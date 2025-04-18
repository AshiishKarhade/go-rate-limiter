package main

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

//var ctx = context.Background()

// TOKEN BUCKET ALGORITHM

const maxTokens = 2       // maximum tokens per user
const refillInterval = 60 // token refills every 60 seconds
const tokensPerRefill = 2

type RateLimiter struct {
	client *redis.Client
}

func (r *RateLimiter) InitializeRateLimit(userID string) {
	_, err := r.client.HGet(ctx, "rate_limit:"+userID, "tokens").Result()
	if errors.Is(err, redis.Nil) {
		r.client.HSet(ctx, "rate_limit:"+userID, "tokens", maxTokens, "last_refill_time", time.Now().Format(time.RFC3339))
	}
}

func (r *RateLimiter) GetTokenBucket(userID string) (int, time.Time, error) {
	tokens, err := r.client.HGet(ctx, "rate_limit:"+userID, "tokens").Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, time.Time{}, fmt.Errorf("count not get the token %v", err)
	}

	lastRefill, err := r.client.HGet(ctx, "rate_limit:"+userID, "last_refill_time").Time()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, time.Time{}, fmt.Errorf("count not get the last refill time %v", err)
	}

	if errors.Is(err, redis.Nil) {
		tokens = maxTokens
		lastRefill = time.Now()
		r.client.HSet(ctx, "rate_limit:"+userID, "tokens", tokens, "last_refill_time", lastRefill)
	}
	return tokens, lastRefill, nil
}

func (r *RateLimiter) RefillTokens(userID string, lastRefill time.Time, tokens int) int {
	elapsed := time.Now().Sub(lastRefill)
	refills := int(elapsed.Seconds()) / refillInterval

	newTokens := tokens + refills*tokensPerRefill
	if newTokens > maxTokens {
		newTokens = maxTokens
	}
	r.client.HSet(ctx, "rate_limit:"+userID, "tokens", newTokens, "last_refill_time", time.Now())
	return newTokens
}

func (r *RateLimiter) AllowRequest(userID string) bool {
	tokens, lastRefill, err := r.GetTokenBucket(userID)
	if err != nil {
		log.Println("Error getting token bucket:", err)
		return false
	}

	// Refill tokens if necessary
	newTokens := r.RefillTokens(userID, lastRefill, tokens)

	if newTokens > 0 {
		// Allow the request and decrease the token count
		r.client.HIncrBy(ctx, "rate_limit:"+userID, "tokens", -1)
		return true
	}

	// Reject the request if there are no tokens available
	return false
}
