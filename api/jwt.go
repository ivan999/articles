package api

import (
    "time"
    "errors"
    "github.com/golang-jwt/jwt"
)

type claimsJWT struct {
    jwt.StandardClaims
    UserID int64
}

const secretKeyJWT = "dsa*#@$*?saafs|fds32*_"

func generateJWT(userID int64) (string, error) {
    const tokenTime = time.Hour * 24

    claims := claimsJWT{
        jwt.StandardClaims{
            ExpiresAt: time.Now().Add(tokenTime).Unix(),
            IssuedAt: time.Now().Unix(),
        },
        userID,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
    tokenString, err := token.SignedString([]byte(secretKeyJWT))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func tokenParser(token *jwt.Token) (any, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
        return []byte(secretKeyJWT), nil
    }

    return nil, errors.New("invalid signing method")
}

func parseJWT(accessToken string) (int64, error) {
    var claims claimsJWT
    _, err := jwt.ParseWithClaims(accessToken, &claims, tokenParser)
    if err != nil {
        return 0, err 
    }

    return claims.UserID, nil
}
