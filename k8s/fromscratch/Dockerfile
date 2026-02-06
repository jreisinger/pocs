FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod fibspin.go ./
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=build /src/fibspin /
ENTRYPOINT ["/fibspin"]
