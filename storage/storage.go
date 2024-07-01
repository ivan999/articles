package storage

import (
    "fmt"
    "database/sql"
    "github.com/go-sql-driver/mysql"
)

type Storage struct {
    db *sql.DB
}

type Credentials struct {
    Username string
    Password string
    IPAddr   string
    DBName   string
}

func Open(cdls *Credentials) (*Storage, error) {
    const funcName = "Open"

    config := mysql.Config{
        User:   cdls.Username,
        Passwd: cdls.Password,
        Addr:   cdls.IPAddr + ":3306",
        Net:    "tcp",
        ClientFoundRows: true,
        AllowNativePasswords: true,
    }

    DSN := config.FormatDSN()

    db, err := sql.Open("mysql", DSN)
    if err != nil {
        return nil, defError(funcName, err)
    }

    tables := []string{
        fmt.Sprintf("DROP DATABASE IF EXISTS %s", cdls.DBName),
        fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cdls.DBName),
        fmt.Sprintf("USE %s", cdls.DBName),
        usersTable, articlesTable,
    }     

    for _, val := range tables {
        _, err := db.Exec(val)
        if err != nil {
            return nil, defError(funcName, err)
        }
    }

    return &Storage{db}, nil
}

func (s *Storage) Close() {
    s.db.Close()
}
