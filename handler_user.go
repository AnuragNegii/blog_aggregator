package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AnuragNegii/blog_aggreagator/internal/database"
	"github.com/google/uuid"
)

func resetTable(s *state, cmd command) error{
	ctx := context.Background()
	if err := s.db.DeleteTable(ctx); err != nil{
		return fmt.Errorf("cant reset table: %v \n", err)
	}

	fmt.Println("Table reset")
	return nil
}

func getAllUsers(s *state, cmd command) error{
	ctx := context.Background()
	allUsers, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("Cant get list of Users: %v", err)
	}

	currentUser := s.cfg.CurrentUserName
	for _, user := range allUsers{
		if user.Name == currentUser {
			fmt.Printf("* %v (current)\n", user.Name)
		}else{
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func handlerLogin(s *state, cmd command) error{
	if len(cmd.Args) < 2{
		return errors.New("no command given")
	}
	name := cmd.Args[1]

	ctx := context.Background()
	user, err := s.db.GetUser(ctx, name); 
	if err != nil {
		return fmt.Errorf("No user found with %v name", name)
	}

	if err := s.cfg.SetUser(user.Name); err != nil{
		return fmt.Errorf("error logging in.")
	}
	fmt.Printf("%v logged in \n", user.Name)
	return nil
}

func registerUser(s *state, cmd command) error{
	if len(cmd.Args) < 2{
		return errors.New("no command given")
	}

	name := cmd.Args[1]
	now := time.Now()

	ctx := context.Background()

	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID: uuid.New(), 
		CreatedAt: now,
		UpdatedAt: now, 
		Name: name,
	})

	if err != nil {
		return errors.New("Couldn't create user") 
	}
	if err := s.cfg.SetUser(user.Name); err != nil {
		return errors.New("Couldn't set username")
	}
	
	fmt.Println("Created user successfully")
	printUser(user)
	return nil
}

func printUser(user database.User){
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)	
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	feed, err := fetchFeed(ctx, url)
	if err != nil{
		return err
	}
	fmt.Printf("%v\n", feed)
	return nil
}
