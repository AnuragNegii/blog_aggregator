package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AnuragNegii/blog_aggreagator/internal/database"
	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/net/html"
)

func resetTable(s *state, cmd command) error{
	ctx := context.Background()
	if err := s.db.DeleteTable(ctx); err != nil{
		return fmt.Errorf("cant reset table: %v \n", err)
	}

	fmt.Println("Tables reset")
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

func handlerAgg(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Not enough arguments\n")
	}
	t, err := time.ParseDuration(cmd.Args[1])
	if err != nil {
		return fmt.Errorf("error parsing time: %v\n", err)
	}
	moreThanThis, err := time.ParseDuration("10s")
	if err != nil {
		return fmt.Errorf("error while setting own time.")
	}
	if t < moreThanThis{
		return fmt.Errorf("time should be at least 10s.")
	}

	ticker := time.NewTicker(t)
	fmt.Printf("Collecting feeds every %v\n", t)
	for ;; <-ticker.C {
		if err  = ScrapeFeed(s, cmd, user); err != nil {
			fmt.Printf("error while scraping feed: %v\n", err)
		}
	}
}

func addFeed(s *state, cmd command, user database.User) error{
	if len(cmd.Args) < 3{
		return fmt.Errorf("Not enough arguments for addFeed command %v, %v, %v", cmd.Args[0], cmd.Args[1], cmd.Args[2])
	}

	ctx := context.Background()
	name := cmd.Args[1]
	providedURL := cmd.Args[2]

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID: uuid.New(),	
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,	
		Url: sql.NullString{String: providedURL, Valid: true},
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("Something went wrong while trying to add feed: %v", err) 
	}
	fmt.Printf("Feed created successfully\n")
	fmt.Printf("Feed Id: %v\n", feed.ID)
	fmt.Printf("Feed Name: %v\n", feed.Name)
	fmt.Printf("Feed url: %v\n", feed.Url)

	insertedFeedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID: uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {

		return fmt.Errorf("error while following: %v", err)
	}

	fmt.Printf("feed Name : %v\n", insertedFeedFollow.FeedName)
	fmt.Printf("Now followed by :%v\n", insertedFeedFollow.UserName)
	
	return nil
}

func handlerFeeds(s *state, cmd command) error{
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving feeds: %v", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Feed: %v\n", feed.Name)
		fmt.Printf("URL: %v\n", feed.Url)
		fmt.Printf("username: %v\n", feed.Username)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error{
	if len(cmd.Args) < 2 {
		return fmt.Errorf("too few arguments\n")
	}
	url := cmd.Args[1]

	ctx := context.Background()
	feed, err := s.db.GetFeed(ctx, sql.NullString{String: url, Valid: true})
	if err != nil {
		return fmt.Errorf("error retrieving feed: %v\n", err)
	}

	insertedFeedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID: uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error while following: %v", err)
	}
	
	fmt.Printf("Followed the feed of this url:%v\n", url)

	fmt.Printf("feed Name : %v\n", insertedFeedFollow.FeedName)
	fmt.Printf("followed by :%v\n", insertedFeedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error{
	ctx := context.Background()
	feedFollow, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error retrieving userFollowList: %v", err)
	}

	fmt.Print("Feeds That user follows")
	for _, followList := range feedFollow {
		fmt.Printf("%v\n", followList.Feedname)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error{
	if len(cmd.Args) < 2 {
		return fmt.Errorf("too few arguments %v\n", len(cmd.Args))
	}
	url := cmd.Args[1]

	ctx := context.Background()
	if err := s.db.Deletefeedfollowrecord(ctx, database.DeletefeedfollowrecordParams{
		UserID: user.ID,
		Url: sql.NullString{String: url, Valid: true},	
	}); err != nil {
		return fmt.Errorf("error while unfollowing: %v\n", err)
	}
	return nil
}

func ScrapeFeed(s *state, cmd command, user database.User)error{
	ctx := context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(ctx, uuid.NullUUID{UUID: user.ID, Valid: true,})
	if err != nil {
		return fmt.Errorf("error getting next feed: %v\n", err)
	}

	if err = s.db.MarkFeedFetched(ctx, nextFeed.ID); err != nil {
		return fmt.Errorf("error marking feed: %v\n", err)
	}

	feed, err := fetchFeed(ctx, nextFeed.Url.String)
	if err != nil{
		return fmt.Errorf("error while fetching feed: %v\n", err)
	}

	for _, item := range feed.Channel.Item{
		t, err := dateparse.ParseAny(item.PubDate)
		if err != nil {
			fmt.Printf("error parsing PubDateL %v\n", err)
			t = time.Now()
		}
		if err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID: uuid.New(),
			Title: item.Title,
			Url: item.Link,
			Description: item.Description,
			PublishedAt: t,
			FeedID: nextFeed.ID,
		}); err != nil{
			var pqError *pq.Error
			if errors.As(err, &pqError){
				if pqError.Code == "23505"{
					continue
				}else{
					return fmt.Errorf("error creating feed: %v\n", err)
				}
			}
		}

	} 
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User)error{
	limit := 2
	if len(cmd.Args) > 1 {
    l, err := strconv.Atoi(cmd.Args[1])
    if err == nil {
		limit = l
    	}
	}
	ctx := context.Background()
	datab, err := s.db.GetPostsForUsers(ctx, database.GetPostsForUsersParams{
		UserID: user.ID,	
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error browsing database: %v\n", err)
	}
	for _, data := range datab{
		fmt.Printf("%v\n",data.Title)
		fmt.Printf("%v\n",htmlToText(data.Description))
		fmt.Printf("%v\n",data.PublishedAt)
		fmt.Println("-----------------------------------------------")
	}

	return nil
}

func middlewareLoggedIn(handler func(s * state, cmd command, user database.User) error) func(*state, command)error {
	return func(s *state, cmd command)error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil{
			return fmt.Errorf("error getting current user %v\n", err)
		}
		return handler(s, cmd, user)
	}
}

func htmlToText(htmlStr string) string {
    doc, err := html.Parse(strings.NewReader(htmlStr))
    if err != nil {
        return htmlStr // fallback if parsing fails
    }
    var f func(*html.Node) string
    f = func(n *html.Node) string {
        if n.Type == html.TextNode {
            return n.Data
        }
        var out string
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            out += f(c)
        }
        return out
    }
    return f(doc)
}
