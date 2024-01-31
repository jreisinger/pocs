FROM golang:1.21 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app

FROM gcr.io/distroless/static-debian12
# FROM gcr.io/distroless/static-debian12:debug
COPY --from=build /bin/app /
CMD ["/app"]
