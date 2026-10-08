package main

import (
	"log"
	"os"

	"github.com/sharipovr/blog-aggregator-go/internal/config"
)

type state struct {
	config *config.Config
}

func main() {

	// Initial read
	var s state
	c, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	s.config = &c

	// p.7 commands
	var cmds commands
	cmds.handlers = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Error: not enough argumanrs provided")
	}

	cmd := command{name: args[1], args: args[2:]}
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
