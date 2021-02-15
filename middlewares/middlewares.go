package middlewares

import "net/http"

type jsonAnswer struct {
	status  string `json:"status"`
	message string `json:"message"`
}

func JsonMiddleware(next http.Handler) {
	
}
