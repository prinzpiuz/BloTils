package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /ping request\n%v", r)
	io.WriteString(w, "Pong\n")
}
