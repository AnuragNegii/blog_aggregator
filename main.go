package main

import (
	"fmt"
	"os"
	"github.com/AnuragNegii/blog_aggregator/internal/config"
)

func main(){
    cfg, err := config.Read() 
    if err != nil{
        fmt.Printf("Error: %s", err)
        return
    }
    configState := state{}   
    configState.config = &cfg
    cmds := commands{}

    cmds.command = make(map[string]func(*state, command) error)
    
    cmds.register("login", handlerLogin)

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
        fmt.Printf("Error: %s", err)
        os.Exit(1)
    }
}
