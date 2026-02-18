package middleware

import (
	"context"
	"net/http"
	"time"

	"go-auth-system/internal/util"

	"github.com/google/uuid"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Create or receive Request ID
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// 2. Set start time
		startTime := time.Now()

		// 3. Inject into Context
		ctx := context.WithValue(r.Context(), util.RequestIDKey, reqID)
		ctx = context.WithValue(ctx, util.StartTimeKey, startTime)

		// 4. Set header to client easy tracking
		w.Header().Set("X-Request-ID", reqID)

		// 5. Continue execute next action/request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
