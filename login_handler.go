package main

import (
	"fmt"
)



func handlerLogin(s *state, cmd command) error{
    if len(cmd.args) == 0{
        return fmt.Errorf("a username is required\n")
    }
    err := s.config.SetUser(cmd.args[0])
    if err != nil{
        return err
    }
    fmt.Printf("user %s has been set\n", cmd.args[0])
    return nil
}
