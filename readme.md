gator : 
1. You need to install Postgres and go to your system
2. clone the repo.
3. You also need to install goose in your folder using this "go install github.com/pressly/goose/v3/cmd/goose@latest"
4. Run `goose postgres $DATABASE_URL up`
5. you can use different commands Like 
    Commands:
    register name - to register a user and login to it in database
    login name - to change user
    addfeed title url - to add an rss feed to ur database which u can follow
    follow url - to follow rss feed
    browse limit - to browse the feeds followed by the user
    agg time(eg - 10s, 15s, 1h) - to get new feeds in the system
