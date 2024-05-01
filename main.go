package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	if len(os.Args) > 1 {
		if os.Args[1] == "--debug" {
			err := os.Remove("db")
			if err != nil {
				log.Println(err)
			}
		}
	}
	err := serverSetup()
	if err != nil {
		fmt.Println(err)
	}
}
