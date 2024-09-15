FROM golang:alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/main .

CMD ["./main"]