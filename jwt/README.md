The aim of jwt is to understand how authorization based on [JWT](https://jwt.io/introduction) tokens works. It contains an API server built with gin. The API has one protected GET endpoint that returns 200 only if correct JWT token is supplied

JWT is an open standard for secure transmission of information as a JSON object. This information can be verified because it's digitally signed. It can be signed:

* using a secret key (HMAC algorithm)
* using a public/private key pair (RSA or ECDSA algorithm)

Although JWTs can be encrypted, this is about signing, i.e. verifying the integrity of the message. When tokens are signed using public/private key pairs, the signature also certifies that only the party holding the private key is the one that signed it.

JWT structure:

* header - contains Base64Url-encoded token type (JWT) and signing algorithm (HMAC, SHA256, RSA)
* payload - contains Base64Url-encoded claims, i.e. statements about an entity ([registered](https://tools.ietf.org/html/rfc7519#section-4.1), public, private)
* signature - created by signing encoded header and payload, e.g. using HMAC algorithm and a secret:

```
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret)
```

Run

```
go run api.go
```

Play

```
# 400: no token in Authorization header
curl -v http://localhost:8080/protected

# 400: token signed with different key than the one defined in `Key` variable
# Using token from https://jwt.io with modified secret key.
curl -v http://localhost:8080/protected -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.sxWNNhU7H0EtNNpwycvDktMxqtLJjR4mHHLwuym0wRs' #gitleaks:allow

# 401: incorrect "sub" claim in the payload
# Using token from https://jwt.io.
curl -v http://localhost:8080/protected -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c' #gitleaks:allow

# 200: correct "sub" claim in the payload
# Using token from https://jwt.io with "sub": "Aquinas"
curl -v http://localhost:8080/protected -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJBcXVpbmFzIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.0_Uqs7hhTP7P_4uTfqQ1AEBvd0cfiRnTEQF2rhtv5aE' #gitleaks:allow
```