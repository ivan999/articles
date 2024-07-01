package api

import (
    "net/http"
    "golang.org/x/crypto/bcrypt"

    "github.com/ivan999/articles/storage"
    "github.com/ivan999/articles/response"
)

func generatePasswordHash(password string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
    return string(hash)
}

func comparePasswordHash(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

const resourceUser = "user"

func (usage *ServerUsage) getUserHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r) 

    userID, err := receiveParam(r, keyID)
    if err != nil {
        handleReceiveError(h, err)
        return
    }

    user, err := usage.Storage.GetUserByID(userID)
    if err != nil {
        details := map[string]any{keyUserID: userID}
        handleStorageError(h, err, details)
        return
    }
    
    user.Password = ""
    h.HandleResponseData(http.StatusOK, "user is found", user)
}

func (usage *ServerUsage) signUpUserHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r)
    details := map[string]any{}

    var user storage.User
    err := receiveJSON(r, &user)
    if err != nil {
        handleReceiveError(h, err)
        return
    }

    user.Password = generatePasswordHash(user.Password)

    userID, err := usage.Storage.AddUser(&user)
    if err != nil {
        details[keyUsername] = user.Username
        handleStorageError(h, err, details)
        return
    }

    details[keyUserID] = userID
    const message = "user is successfuly created"
    h.HandleResponseDetails(http.StatusCreated, message, details)
}

func doVerifying(
    h *response.ResponseHandler, userID int64, password, hash string,
) {
    details := map[string]any{}
    if comparePasswordHash(password, hash) {
        token, err := generateJWT(userID)
        if err != nil {
            h.HandleServerError("failed to generate token", err)
            return
        }
        details[keyToken] = token
        details[keyUserID] = userID
        h.HandleResponseDetails(http.StatusOK, "user is verified", details)
    } else {
        details[keyPassword] = password
        const message = "invalid password"
        h.HandleClientError(http.StatusUnauthorized, message, details)
    }
}

func (usage *ServerUsage) signInUserHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r)

    var user storage.User
    err := receiveJSON(r, &user)
    if err != nil {
        handleReceiveError(h, err)
        return 
    }

    storedUser, err := usage.Storage.GetUserByName(user.Username)
    if err != nil {
        details := map[string]any{keyUsername: user.Username}
        handleStorageError(h, err, details)
        return
    }

    doVerifying(h, storedUser.UserID, user.Password, storedUser.Password)
}

func (usage *ServerUsage) updateUserHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r)
    details := map[string]any{}

    userID := r.Context().Value(keyUserID).(int64)
    
    var user storage.User
    err := receiveJSON(r, &user)
    if err != nil {
        handleReceiveError(h, err)
        return
    }

    user.Password = generatePasswordHash(user.Password)

    details[keyUserID] = userID
    err = usage.Storage.UpdateUser(userID, &user)
    if err != nil {
        handleStorageError(h, err, details)
        return
    }
    
    const message = "user is successfuly updated"
    h.HandleResponseDetails(http.StatusOK, message, details)
}

func (usage *ServerUsage) deleteUserHandler(
    w http.ResponseWriter, r *http.Request,
) {
    h := response.NewResponseHandler(w, r)
    details := map[string]any{}

    userID := r.Context().Value(keyUserID).(int64)

    details[keyUserID] = userID 
    err := usage.Storage.DeleteUser(userID)
    if err != nil {
        handleStorageError(h, err, details)
        return
    }

    const message = "user is successfuly deleted"
    h.HandleResponseDetails(http.StatusOK, message, details)
}
