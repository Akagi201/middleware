# jwt middleware

A middleware that will check that a [JWT](http://jwt.io/) is sent on the `Authorization` header and will then set the content of the JWT into the `user` variable of the request.

This module lets you authenticate HTTP requests using JWT tokens in your Go Programming Language applications. JWTs are typically used to protect API endpoints, and are often issued using OpenID Connect.

## Key Features

* Ability to **check the `Authorization` header for a JWT**
* **Decode the JWT** and set the content of it to the request context

### Token Extraction

The default value for the `Extractor` option is the `FromAuthHeader`
function which assumes that the JWT will be provided as a bearer token
in an `Authorization` header, i.e.,

```
Authorization: bearer {token}
```

To extract the token from a query string parameter, you can use the
`FromParameter` function, e.g.,

```go
jwtmiddleware.New(jwtmiddleware.Options{
  Extractor: jwtmiddleware.FromParameter("auth_code"),
})
```

In this case, the `FromParameter` function will look for a JWT in the
`auth_code` query parameter.

Or, if you want to allow both, you can use the `FromFirst` function to
try and extract the token first in one way and then in one or more
other ways, e.g.,

```go
jwtmiddleware.New(jwtmiddleware.Options{
  Extractor: jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader,
                                     jwtmiddleware.FromParameter("auth_code")),
})
```
