FROM golang:1.20 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app

FROM gcr.io/distroless/static-debian11
# FROM gcr.io/distroless/static-debian11:debug
COPY --from=build /bin/app /
CMD ["/app"]