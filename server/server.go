package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
	// login sucress
	// fmt.Println("1.) cfc-Create Flashcard")
	// var userinput string
	// fmt.Scanln(&userinput)
	// if userinput == "cfc" {
	// 	create.Createfc()
	// }
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
