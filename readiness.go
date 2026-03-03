package main

import (
	"net/http"
	
)


func handlerReadiness(w http.ResponseWriter, r *http.Request) {                      // For simplicity, we assume the service is always ready. In a real application, you would check dependencies here.
	w.Header().Add("Content-Type", "text/plain; charset=utf-8\r\n")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}