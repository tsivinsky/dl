package commander

import (
	"flag"
	"fmt"
	"os"
)

type CommandFunc func() error

type Command struct {
	description string
	function    CommandFunc
}

type commander struct {
	commands map[string]Command
}

func (c *commander) RegisterCommand(name string, description string, function CommandFunc) {
	c.commands[name] = Command{
		description: description,
		function:    function,
	}
}

func (c *commander) AddAliases(commandName string, aliases ...string) {
	for _, alias := range aliases {
		c.commands[alias] = c.commands[commandName]
	}
}

func (c *commander) RunCommand(name string) error {
	command, exists := c.commands[name]
	if !exists {
		return fmt.Errorf("command %s not registered", name)
	}

	return command.function()
}

func (c *commander) Usage() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	for name, command := range c.commands {
		fmt.Printf("  %s - %s\n", name, command.description)
	}
	flag.PrintDefaults()
}

func New() *commander {
	return &commander{
		commands: make(map[string]Command),
	}
}
