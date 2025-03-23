package main

import "fmt"

// import "fmt"



type commands struct{
    command map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
    if c == nil{
        fmt.Println("Command struct is not initialized")
        return
    }
    if c.command == nil {
        c.command = make(map[string]func(*state, command) error)
    }
    c.command[name] = f
}

func (c *commands) run (s *state, cmd command) error{
    command, exists := c.command[cmd.name]
    if exists{
        err :=  command(s, cmd) 
        return err
    }
    return fmt.Errorf("the command does not exist.") 
}
