FROM golang:1.22.3-alpine AS build
RUN apk add --no-cache gcc libc-dev
WORKDIR /go/src/app

COPY . .
RUN go test  ./...
RUN go build -o /bin/mutating-wh


FROM alpine:3.20.0
RUN apk add --no-cache ca-certificates

COPY --from=build /bin/mutating-wh /usr/local/bin/mutating-wh
CMD ["mutating-wh"]
