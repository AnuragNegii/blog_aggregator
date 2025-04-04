package main

import (
	"database/sql"
	"fmt"
	"os"
	"github.com/AnuragNegii/blog_aggregator/internal/config"
	"github.com/AnuragNegii/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main(){
    cfg, err := config.Read() 
    if err != nil{
        fmt.Printf("Error: %s", err)
        return
    }
    configState := state{}   
    configState.config = &cfg
    db, err := sql.Open("postgres", cfg.DbURL)
    if err != nil{
        fmt.Printf("Error: %s\n", err)
        os.Exit(1)
    }
    dbQueries := database.New(db) 
    configState.db = dbQueries
    cmds := commands{}
    cmds.command = make(map[string]func(*state, command) error)
    cmds.register("login", handlerLogin)
    cmds.register("register", handlerRegister)
    cmds.register("reset", handlerReset)
    cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerListFeeds)
    //Command line arguments passed by the user
    args := os.Args
    if len(args) < 2 {
        fmt.Print("not enough arguments were provided\n")
        os.Exit(1) 
    }
    cmdName := args[1]
    cmdArgs:= args[2:]
    cmd := command{
        name: cmdName,
        args: cmdArgs,
    }
    err = cmds.run(&configState, cmd)
    if err != nil{
        fmt.Printf("Error: %s\n", err)
        os.Exit(1)
    }
}
