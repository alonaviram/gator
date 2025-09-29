package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/alonaviram/gator/internal/config"
	"github.com/alonaviram/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
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

	db, err := sql.Open("postgres", config.DBURL)
	if err != nil {
		panic(fmt.Sprintf("Couldn't connect to db url:%v", config.DBURL))
	}
	databaseQueries := database.New(db)
	state := state{
		db:     databaseQueries,
		config: &config,
	}

	commands := commands{
		cmds: createCommandsMap(),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerGetUsers)

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
