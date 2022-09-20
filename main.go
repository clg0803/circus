package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/clg0803/circus/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey lang! 🐵🙊🙉🙈\n", user.Username)
	fmt.Printf("Feel free to type in commands \n")
	repl.Start(os.Stdin, os.Stdout)
}
