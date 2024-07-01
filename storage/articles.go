package storage

import (
    "fmt"
    "database/sql"
)

const articlesTable = `
CREATE TABLE IF NOT EXISTS articles (
    article_id INT NOT NULL AUTO_INCREMENT,
    user_id    INT NOT NULL,
    title      TEXT(2048) NOT NULL,
    content    TEXT NOT NULL,

    PRIMARY KEY (article_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
)`

type Article struct {
    ArticleID int64  `json:"articleID"`
    UserID    int64  `json:"userID"`
    Title     string `json:"title"`
    Content   string `json:"content"`
}

const resourceArticle = "article"

func (s *Storage) AddArticle(userID int64, article *Article) (int64, error) {
    const funcName = "AddArticle"

    result, err := s.db.Exec(`
        INSERT INTO articles (user_id, title, content)
        VALUES (?, ?, ?)`,
        userID, article.Title, article.Content,
    )
    if err != nil {
        return 0, defArticleExecError(funcName, err)
    }

    articleID, err := result.LastInsertId()
    if err != nil {
        return 0, defError(funcName, err)
    }

    return articleID, nil
}

func (s *Storage) UpdateArticle(
    articleID, userID int64, article *Article,
) error {
    const funcName = "UpdateArticle"
    
    result, err := s.db.Exec(`
        UPDATE articles SET title = ?, content = ?
        WHERE article_id = ? AND user_id = ?`,
        article.Title, article.Content,
        articleID, userID,
    )
    if err != nil {
        return defArticleExecError(funcName, err)
    }

    return defUpdateError(funcName, resourceArticle, result)
}

func (s *Storage) DeleteArticle(articleID, userID int64) error {
    const funcName = "DeleteArticle"

    result, err := s.db.Exec(
        "DELETE FROM articles WHERE article_id = ? AND user_id = ?",
        articleID, userID,
    )
    if err != nil {
        return defError(funcName, err)
    }

    return defDeleteError(funcName, resourceArticle, result)
}

func articlesHeadersFromRows(rows *sql.Rows, capacity int64) ([]Article, error) {
    headers := make([]Article, 0, capacity)
    for rows.Next() {
        var header Article
        err := rows.Scan(&header.ArticleID, &header.UserID, &header.Title)
        if err != nil {
            return nil, err
        }
        headers = append(headers, header)
    }

    err := rows.Err()
    if err != nil {
        return nil, err
    }

    return headers, nil
}

func (s *Storage) GetArticlesHeaders(offset, limit int64) ([]Article, error) {
    const funcName = "GetArticlesHeaders"
    
    query := "SELECT article_id, user_id, title FROM articles"
    query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, defQueryError(funcName, resourceArticle, err)
    }
    defer rows.Close()

    headers, err := articlesHeadersFromRows(rows, limit)
    if err != nil {
        return nil, defError(funcName, err)
    }
    
    return headers, nil
}

func (s *Storage) GetUserArticlesHeaders(userID int64) ([]Article, error) {
    const (
        funcName = "GetUserArticlesHeaders"
        capacity = 10
    )

    rows, err := s.db.Query(`
        SELECT article_id, user_id, title FROM articles 
        WHERE user_id = ?`,
        userID,
    )
    if err != nil {
        return nil, defQueryError(funcName, resourceArticle, err)
    }
    defer rows.Close()

    headers, err := articlesHeadersFromRows(rows, capacity)
    if err != nil {
        return nil, defError(funcName, err)        
    }
    
    return headers, nil
}

func (s *Storage) GetArticle(articleID int64) (*Article, error) {
    const funcName = "GetArticle"

    row := s.db.QueryRow(
        "SELECT * FROM articles WHERE article_id = ?",
        articleID,
    )

    var article Article
    err := row.Scan(
        &article.ArticleID, &article.UserID, 
        &article.Title, &article.Content,
    )
    if err != nil {
        return nil, defQueryError(funcName, resourceArticle, err)
    }

    return &article, nil
}
