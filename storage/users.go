package storage

import (
    "fmt"
    "database/sql"
)

const usersTable = `
CREATE TABLE IF NOT EXISTS users (
    user_id    INT NOT NULL AUTO_INCREMENT,
    first_name CHAR(64) NOT NULL,
    last_name  CHAR(64) NOT NULL,
    username   CHAR(64) UNIQUE NOT NULL,
    password   CHAR(255) NOT NULL,

    PRIMARY KEY (user_id)
)`

type User struct {
    UserID    int64  `json:"userID"`
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Username  string `json:"username"`
    Password  string `json:"password"`
}

const resourceUser = "user"

func (s *Storage) AddUser(user *User) (int64, error) {
    const funcName = "AddUser"

    result, err := s.db.Exec(`
        INSERT INTO users (first_name, last_name, username, password)
        VALUES (?, ?, ?, ?)`,
        user.FirstName, user.LastName, 
        user.Username, user.Password,
    )
    if err != nil {
        return 0, defUserExecError(funcName, err)
    }
    
    userID, err := result.LastInsertId()
    if err != nil {
        return 0, defError(funcName, err)
    }

    return userID, nil
}

func (s *Storage) UpdateUser(userID int64, user *User) error {
    const funcName = "UpdateUser"

    result, err := s.db.Exec(`
        UPDATE users 
        SET first_name = ?, last_name = ?, username = ?, password = ?
        WHERE user_id = ?`,
        user.FirstName, user.LastName, user.Username,
        user.Password, userID,
    )
    if err != nil {
        return defUserExecError(funcName, err)
    }

    return defUpdateError(funcName, resourceUser, result)
}

func (s *Storage) DeleteUser(userID int64) error {
    const funcName = "DeleteUser"

    result, err := s.db.Exec("DELETE FROM users WHERE user_id = ?", userID)
    if err != nil {
        return defError(funcName, err)
    }

    return defDeleteError(funcName, resourceUser, result)
}

func getUserFromRow(row *sql.Row) (*User, error) {
    var user User
    err := row.Scan(
        &user.UserID, &user.FirstName, &user.LastName, 
        &user.Username, &user.Password,
    )
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (s *Storage) getUser(key string, value any) (*User, error) {
    const funcName = "GetUserByName"

    query := fmt.Sprintf("SELECT * FROM users WHERE %s = ?", key)
    row := s.db.QueryRow(query, value)
    user, err := getUserFromRow(row)
    if err != nil {
        return nil, defQueryError(funcName, resourceUser, err)
    }

    return user, nil
}

func (s *Storage) GetUserByID(userID int64) (*User, error) {
    return s.getUser("user_id", userID)
}

func (s *Storage) GetUserByName(username string) (*User, error) {
    return s.getUser("username", username)
}
