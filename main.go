package main

import (
	"fmt"
	"strings"

	"github.com/popoiuioopp/Learning-hub-Backend/server"
)

func EqualFold(s, t string) bool {
	return strings.EqualFold(s, t)
}

func login() {
	fmt.Println("Enter Code or login or register: ")

	// var then variable name then variable type
	var first string

	// Taking input from user
	fmt.Scanln(&first)
	if EqualFold(first, "login") {
		fmt.Println("Enter username : ")
		var username string
		fmt.Scanln(&username)

		fmt.Println("Enter password : ")
		var password string
		fmt.Scanln(&password)

		// if login sucess

	} else if EqualFold(first, "register") {
		fmt.Println(">>>>>>>Register")
		fmt.Println("Username : ")
		var regisUsername string
		fmt.Scanln(&regisUsername)

		fmt.Println(">>>>>>>Register")
		fmt.Println("Password : ")
		var regisPass string
		fmt.Scanln(&regisPass)

		// fmt.Println(">>>>>>>Register")
		// fmt.Println("Verify Password :")
		// var regisVerifyPass string
		// fmt.Scanln(&regisVerifyPass)

		// for regisPass != regisPass {
		// 	fmt.Println(">>>>>>>Register : Password do now match")
		// 	fmt.Println("Verify Password :")
		// 	var regisVerifyPass string
		// 	fmt.Scanln(&regisVerifyPass)
		// }
	} else {
		fmt.Println("http://localhost/" + first)
	}
}
func main() {
	/*
		// fmt.Println("Hello World! Tangy Ma Laew")
		// Println function is used to
		// display output in the next line
		for {
			login()
			if false {
				break
			}
		}
		var wait string
		fmt.Scanln(&wait)
	*/
	server.Server()
}
