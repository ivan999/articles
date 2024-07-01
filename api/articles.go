package api

import (
    "net/http"

    "github.com/ivan999/articles/storage"
    "github.com/ivan999/articles/response"
)

const resourceArticle = "article"

func (usage *ServerUsage) createArticleHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r)
    details := map[string]any{}

    userID := r.Context().Value(keyUserID).(int64)
    
    var article storage.Article
    err := receiveJSON(r, &article)
    if err != nil {
        handleReceiveError(h, err)
        return
    }

    details[keyUserID] = userID
    articleID, err := usage.Storage.AddArticle(userID, &article)
    if err != nil {
        handleStorageError(h, err, details)
        return
    }
    
    details[keyArticleID] = articleID
    const message = "article is successfuly created"
    h.HandleResponseDetails(http.StatusCreated, message, details)
} 

func (usage *ServerUsage) updateArticleHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r)
    details := map[string]any{}

    userID := r.Context().Value(keyUserID).(int64)

    articleID, err := receiveParam(r, keyID)
    if err != nil {
        handleReceiveError(h, err)
        return
    }

    var article storage.Article
    err = receiveJSON(r, &article)
    if err != nil {
        handleReceiveError(h, err)
        return
    }

    details[keyUserID] = userID
    details[keyArticleID] = articleID
    err = usage.Storage.UpdateArticle(articleID, userID, &article)
    if err != nil {
        handleStorageError(h, err, details)
        return
    }

    const message = "article is successfuly updated"
    h.HandleResponseDetails(http.StatusOK, message, details)
}

func (usage *ServerUsage) deleteArticleHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r)
    details := map[string]any{}

    userID := r.Context().Value(keyUserID).(int64)

    articleID, err := receiveParam(r, keyID)
    if err != nil {
        handleReceiveError(h, err)
        return
    }

    details[keyUserID] = userID
    details[keyArticleID] = articleID
    err = usage.Storage.DeleteArticle(articleID, userID)
    if err != nil {
        handleStorageError(h, err, details)
        return
    }

    const message = "article is successfuly deleted"
    h.HandleResponseDetails(http.StatusOK, message, details)
}

func (usage *ServerUsage) getArticlesHeadersHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r)
    details := map[string]any{}

    offset, err := receiveParam(r, keyOffset)
    if err != nil {
        handleReceiveError(h, err)
        return
    }
    limit, err := receiveParam(r, keyLimit)
    if err != nil {
        handleReceiveError(h, err)
        return
    }

    details[keyLimit] = limit
    details[keyOffset] = offset
    headers, err := usage.Storage.GetArticlesHeaders(offset, limit)
    if err != nil {
        handleStorageError(h, err, details)
        return
    }

    h.HandleResponseData(http.StatusOK, "headers are found", headers)
}

func (usage *ServerUsage) getArticleHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r)

    articleID, err := receiveParam(r, keyID)
    if err != nil {
        handleReceiveError(h, err)
        return
    }

    article, err := usage.Storage.GetArticle(articleID)
    if err != nil {
        details := map[string]any{keyArticleID: articleID}
        handleStorageError(h, err, details)
        return
    }

    h.HandleResponseData(http.StatusOK, "article is found", article)
}

func (usage *ServerUsage) getUserArticlesHeadersHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r)

    userID, err := receiveParam(r, keyUserID)
    if err != nil {
        handleReceiveError(h, err)
        return
    }

    headers, err := usage.Storage.GetUserArticlesHeaders(userID)
    if err != nil {
        details := map[string]any{keyUserID: userID}
        handleStorageError(h, err, details)
    }

    h.HandleResponseData(http.StatusOK, "headers are found", headers)
}
