package main

import (
	"context"
	"fmt"
	"time"
	"github.com/AnuragNegii/blog_aggregator/internal/database"
	"github.com/google/uuid"
)


func handlerLogin(s *state, cmd command) error{
    if len(cmd.args) == 0{
        return fmt.Errorf("a username is required\n")
    }
    userName := cmd.args[0]
    context := context.Background()
    _,err := s.db.GetUser(context, userName)
    if err != nil{
        return fmt.Errorf("No user found")
    }
    err = s.config.SetUser(userName)
    if err != nil{
        return err
    }
    fmt.Printf("user %s has been set\n", cmd.args[0])
    return nil
}


func handlerRegister(s *state, cmd command) error{
    if len(cmd.args) == 0{
        return fmt.Errorf("no username provided") 
    }
    username := cmd.args[0]
    context := context.Background()
    _, err := s.db.GetUser(context, username)
    if err == nil{
        return fmt.Errorf("User already exist: %s", username)
    }
    newUser, err := s.db.CreateUser(context, database.CreateUserParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: username,
    })
    if err != nil {
        return fmt.Errorf("Error: %w", err)
    }
    err = s.config.SetUser(newUser.Name)
    if err != nil{
        return fmt.Errorf("error setting user: %w", err)
    }
    fmt.Printf("New user created: %s\n", username)
    return nil
}
