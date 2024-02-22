FROM golang:1.22 as build
WORKDIR /fibspin
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -o /bin/fibspin

FROM scratch
COPY --from=build /bin/fibspin /
ENTRYPOINT ["/fibspin"]
