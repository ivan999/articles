# Articles

REST API for creating and reading articles. 
Works with JSON. 
Requires authorization with JWT tokens.

## The user structure:
    User ID
    First name
    Last name
    Username
    Password
  
## The user endpoints:
  ### /user/sign-up
    If there isn't a user with the given username, creates a new user.
  
  ### /user/sign-in
    If there is a user with the given username and the given password is correct, gets the user ID and a JWT token.
  
  ### /user/update
    Parses the JWT token from the HTTP header, gets the user ID from it, and update the user.
  
  ### /user/delete
    Parses the JWT token from the HTTP header, gets the user ID from it, and delete the user.
  
  ### /user?id=?
    Gives the user by the user ID given in the query parameter.
  

## The article structure:
    Article ID
    User ID
    Title
    Content

## The article endpoints:
  ### /article/create
    Parses the JWT token from the HTTP header, creates a new article.

  ### /article/update?id=?
    Parses the JWT token from the HTTP header, updates the article.

  ### /article/delete?id=?
    Parses the JWT token from the HTTP header, deletes the article.

  ### /article?id=?
    Gives the article by the article ID given in the query parameter.

  ### /articles/headers?limit=?&offset=?
    Gives article headers (part of the article structure without the content field) based on the limit and offset query parameters.

  ### /articles/headers?userID=?
    Gives article headers based on the user ID query parameter.
