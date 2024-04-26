package main

import "fmt"

func main() {
	err := serverSetup()
	if err != nil {
		fmt.Println(err)
	}
}
