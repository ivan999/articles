package main

import (
    "log"
    "github.com/ivan999/articles/storage"
    "github.com/ivan999/articles/api"
)

const (
    logFilePath = "app.log"
    serverPort = "8080"
)

var credentials = storage.Credentials{
    Username: "public",
    Password: "",
    IPAddr:   "127.0.0.1",
    DBName:   "test_articles",
}

func main() {
    storage, err := storage.Open(&credentials)
    if err != nil {
        log.Fatal(err) 
    }

    usage := api.ServerUsage{Storage: storage}
    log.Fatal(api.RunServer(serverPort, &usage))
}
