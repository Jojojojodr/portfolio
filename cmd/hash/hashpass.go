package main

import (
	"fmt"
	"os"

	"github.com/Jojojojodr/portfolio/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hashpass <password>")
		return
	}
	password := os.Args[1]
	hash := internal.Encrypt(password)
	fmt.Println(hash)
}