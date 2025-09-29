package main

import (
	"fmt"
	"os"

	"github.com/alonaviram/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}
type commands struct {
	cmds map[string]func(*state, command) error
}

func (c commands) run(s *state, cmd command) error {
	callback, exists := c.cmds[cmd.name]
	if !exists {
		return fmt.Errorf("command %v is not exists", cmd.name)
	}
	err := callback(s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func main() {
	config, err := config.Read()
	if err != nil {
		panic("wtf")
	}

	state := state{
		config: &config,
	}

	commands := commands{
		cmds: createCommandsMap(),
	}
	commands.register("login", handlerLogin)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("no argument received")
		os.Exit(1)
	}
	err = commands.run(&state, command{
		name: args[1],
		args: args[2:],
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func createCommandsMap() map[string]func(*state, command) error {
	commandsMap := make(map[string]func(*state, command) error)
	return commandsMap
}
