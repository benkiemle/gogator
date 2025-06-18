package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	s, err := NewState()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer s.connection.Close()
	cmds := GetCommands()

	args := os.Args

	if len(args) < 2 {
		fmt.Println("ERROR: Not enough arguments provided")
		os.Exit(1)
	}

	commandName := args[1]
	cmd := command{
		name: commandName,
		args: args[2:],
	}
	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
