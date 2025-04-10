# STAGE 1: Build
FROM golang:1.22 AS build
WORKDIR /src/wizards/
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/wizards

# STAGE 2: Deployment
FROM gcr.io/distroless/static-debian12
USER nobody:nobody
COPY --from=build /bin/wizards /
CMD [ "/wizards" ]
