# Intro to JWT authentication middleware with chi

```sh
$ go run main.go
$ curl -X POST localhost:3000/login
$ curl -H "Authorization: Bearer <your_token>" localhost:3000/protected
```

It's generally not a good practice to generate JWT tokens in the same web server that consumes them. Here's why:

- Coupling: authn logic is tightly coupled with your app logic
- Scalability: hard to share authn across multiple services
- Security (SPoF): JWT secret is exposed in every service that need to verify tokens

When the same-service generation might be OK:

- Small apps with single service
- PoC or prototypes
- Internal tools with limited scope

Better architecture patterns:

(1) Separate authn service

```go
// Auth Service (separate microservice)
func authService() {
    r := chi.NewRouter()
    r.Post("/login", loginHandler)
    r.Post("/refresh", refreshHandler)
    // Only this service generates tokens
}

// Your Application Service
func appService() {
    r := chi.NewRouter()
    r.With(jwtMiddleware).Get("/protected", protectedHandler)
    // Only verifies tokens, doesn't generate them
}
```

(2) Use an external identity provider (IDP)

- Auth0, Okta, Firebase Auth
- OAuth2/OpenID Connect providers (Google, Github, etc.)
- AWS Cognito, Azure AD B2C

(3) API gateway

Handle authn at the gateway level. Your services receive only verified requests.

- Kong, Traefik, or AWS API Gateway