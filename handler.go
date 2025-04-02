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

func handlerReset(s *state, cmd command) error{
    ctx := context.Background()
    err := s.db.ResetUser(ctx)
    if err != nil{
        return fmt.Errorf("Error: %s", err)
    }
    fmt.Println("Database reset done....!")
    return nil
}

func handlerGetUsers(s *state, cmd command) error{
    ctx := context.Background()
    users, err := s.db.GetUsers(ctx)
    if err != nil{
        return fmt.Errorf("Error: %s",err)
    }
    for user := range users{
        if users[user] == s.config.CurrentUserName{
            fmt.Printf("%s (current)\n", users[user])
            continue
        }
        fmt.Println(users[user])
    }
    return nil
}

func handlerAgg(s *state, cmd command) error{
	ctx := context.Background()
	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err 
	}
	fmt.Println(feed.Channel.Item[0].Title)
	fmt.Println(feed.Channel.Item[0].Description)


    return nil
}


func handlerAddFeed(s *state, cmd command) error{
	ctx := context.Background()
	if len(cmd.args)< 2{
		return fmt.Errorf("usage: addfeed <name> url <url>")
	}
	feedName  := cmd.args[0]
	feedURL := cmd.args[1]
	currentUser, err  :=  s.db.GetUser(ctx, s.config.CurrentUserName)
	if err != nil{
		return err
	}
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: feedName,
		Url: feedURL,
		UserID: currentUser.ID,
	})
	if err != nil{
		return fmt.Errorf("Could not create feed: %v", err) 
	}
 	// Print the feed details
    fmt.Printf("Feed created: %+v\n", feed)
    return nil
}

func handlerListFeeds(s *state, cmd command) error{
	feed, err := s.db.ListFeeds(context.Background())	
	if err != nil {
		return err
	}
	for f := range feed{
		fmt.Printf("feedName: %s\n", feed[f].Name)
		fmt.Printf("URL: %s\n", feed[f].Url)
		fmt.Printf("UserName: %s\n", feed[f].UserName)
	}
	return nil
}
