FROM golang:1.22 as build
WORKDIR /wizards
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/wizards

FROM gcr.io/distroless/static-debian12
COPY --from=build /bin/wizards /
CMD ["/wizards"]
