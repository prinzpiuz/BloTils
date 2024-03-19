package handlers

import (
	// "fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func GetClaps(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	url := params["url"]
	print(url)
	println("")
	// fmt.Printf("got / request\n%v", r)
	io.WriteString(w, "This is my website!\n")
}
