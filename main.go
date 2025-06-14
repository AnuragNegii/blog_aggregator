package main

import (
	"database/sql"
	"log"
	"os"

	config "github.com/AnuragNegii/blog_aggreagator/internal/config"
	"github.com/AnuragNegii/blog_aggreagator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}

func main(){
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("err reading config file")
	}
	
	newState := state{}
	newState.cfg = &cfg

	cmds := commands{
		Handler: make(map[string]func(*state, command) error),
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil{
		log.Fatalf("Cant connect to db")
	}
	defer db.Close()

	dbQuerries := database.New(db)
	newState.db = dbQuerries 

	cmds.register("login", handlerLogin)
	cmds.register("register", registerUser)
	cmds.register("reset", resetTable)
	cmds.register("users", getAllUsers)
	cmds.register("agg", middlewareLoggedIn(handlerAgg))
	cmds.register("addfeed", middlewareLoggedIn(addFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	if len(os.Args) < 2 {
		log.Fatal("no command given")
	}

	userCmd := command{}
	userCmd.name =os.Args[1] // name of the command
	userCmd.Args = os.Args[1:] // rest of the arguments

	if err := cmds.run(&newState, userCmd); err != nil {
		log.Fatalf("error running command: %v\n", err)
	}

}
