FROM golang:1.22 as build
WORKDIR /src
COPY go.mod fibspin.go ./
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=build /src/fibspin /
ENTRYPOINT ["/fibspin"]
