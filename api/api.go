package api

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"

    "github.com/ivan999/articles/storage"
)

type ServerUsage struct {
    Storage *storage.Storage
}

const (
    keyID = "id"
    expID = "{id:[0-9]+}"

    keyOffset = "offset"
    keyLimit = "limit"

    expOffset = "{offset:[0-9]+}"
    expLimit = "{limit:[0-9]+}"

    keyUserID = "userID"
    keyArticleID = "articleID"

    expUserID = "{userID:[0-9]+}"
    expArticleID = "{articleID:[0-9]+}"

    keyError = "error"
    keyToken = "token"
    keyHeader = "header"
    keyUsername = "username"
    keyPassword = "password"
)

func RunServer(port string, usage *ServerUsage) error {
    router := mux.NewRouter()

    usersRouting(router, usage)
    articlesRouting(router, usage)
    
    return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

func usersRouting(router *mux.Router, usage *ServerUsage) {
    router.HandleFunc("/user/sign-up", usage.signUpUserHandler).Methods("POST")

    router.HandleFunc("/user/sign-in", usage.signInUserHandler).Methods("POST")

    router.HandleFunc("/user/update", 
        usage.authHandler(usage.updateUserHandler)).Methods("PUT")

    router.HandleFunc("/user/delete", 
        usage.authHandler(usage.deleteUserHandler)).Methods("POST")

    route := router.HandleFunc("/user", usage.getUserHandler).Methods("GET")
    route.Queries(keyID, expID)
}

func articlesRouting(router *mux.Router, usage *ServerUsage) {
    router.HandleFunc("/article/create", 
        usage.authHandler(usage.createArticleHandler)).Methods("POST")

    route := router.HandleFunc("/article/update",
        usage.authHandler(usage.updateArticleHandler)).Methods("PUT")
    route.Queries(keyID, expID) 

    route = router.HandleFunc("/article/delete",
        usage.authHandler(usage.deleteArticleHandler)).Methods("POST")
    route.Queries(keyID, expID)

    route = router.HandleFunc("/articles/headers", 
        usage.getArticlesHeadersHandler).Methods("GET")
    route.Queries(keyOffset, expOffset, keyLimit, expLimit)

    route = router.HandleFunc("/articles/headers", 
        usage.getUserArticlesHeadersHandler).Methods("GET")
    route.Queries(keyUserID, expUserID)

    route = router.HandleFunc("/article", 
        usage.getArticleHandler).Methods("GET")
    route.Queries(keyID, expID)
}
