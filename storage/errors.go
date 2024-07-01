package storage

import (
    "fmt"
    "database/sql"
    "github.com/go-sql-driver/mysql"
)

const (
    ErrCodeNotFound = 1
    ErrCodeUniqueExists = 2
)

type StorageError struct {
    Code     int
    Message  string
    Wrapped  error
}

func (serr *StorageError) Error() string {
    return fmt.Sprintf(
        "storage error code=%v message='%s' %v",
        serr.Code, serr.Message, serr.Wrapped,
    )
}

func wrapError(funcName string, err error) error {
    return fmt.Errorf("function=%s value=%w", funcName, err)
}

func defError(funcName string, err error) error {
    return &StorageError{
        Code: 0,
        Message: "storage error",
        Wrapped: wrapError(funcName, err),
    }
}

func defNotFound(funcName, resourceName string, err error) error {
    return &StorageError{
        Code: ErrCodeNotFound,
        Message: resourceName + " is not found",
        Wrapped: wrapError(funcName, err),
    }
}

func defUniqueExists(funcName, resourceName string, err error) error {
    return &StorageError{
        Code: ErrCodeUniqueExists,
        Message: resourceName + " exists",
        Wrapped: wrapError(funcName, err),
    }
}

func defUserExecError(funcName string, err error) error {
    if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == 1062 {
        return defUniqueExists(funcName, "username", err)
    }
    return defError(funcName, err)
}

func defExecError(funcName, resourceName string, err error) error {
    if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == 1452 {
        return defNotFound(funcName, resourceName, err)
    }
    return defError(funcName, err)
}

func defArticleExecError(funcName string, err error) error {
    return defExecError(funcName, resourceUser, err)
}

func defCommentExecError(funcName string, err error) error {
    return defExecError(funcName, resourceUser + " or " + resourceArticle, err)
}

func defQueryError(funcName, resourceName string, err error) error {
    if err == sql.ErrNoRows {
        return defNotFound(funcName, resourceName, err)
    }
    return defError(funcName, err)
}

func defUpdateError(funcName, resourceName string, result sql.Result) error {
    affected, err := result.RowsAffected()
    if err != nil {
        return defError(funcName, err) 
    }

    if affected == 0 {
        return defNotFound(funcName, resourceName, nil)
    }

    return nil
}

func defDeleteError(funcName, resourceName string, result sql.Result) error {
    return defUpdateError(funcName, resourceName, result)
}
