tokenauth demonstrates how token-based authentication work. This is the lifecycle of a token:

1. User authenticates using username and password.
2. Credentials are verified agains a database or authentication provider (like AWS Cognito).
3. A token is created and digitally signed to protect its integrity.
4. Token is sent back to the user.
5. Client stores the token and includes it in subsequent requests.
6. The service verifies to token's integrity by validating its signature. If valid it extracts the user's identity and other details (claims) from the token and processes the request.
7. Once the token is verified the user can access the service until the token expires or is invalidated.

JSON Web Token (JWT) is used. 

```
go run ./cmd/register jack
Enter password: jack123

go run ./cmd/server
curl localhost:8080/authenticate -d "username=jack" -d "password=jack123"   # 1. - 4.
curl localhost:8080/verify -H "Authorization: Bearer TOKEN"                 # 5. - 7.
```

For more see Cloud Native Go book -> 12. Security -> Authentication -> Token-based Authentication.