package main

import (
	"fmt"
	"net/http"

	"github.com/popoiuioopp/Learning-hub-Backend/server/create"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
	create.Createfc()
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
