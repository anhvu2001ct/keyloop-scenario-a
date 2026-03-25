FROM golang:1.26.1-alpine3.23 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/sql-migrate
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/appointment-service

FROM alpine:3.23
RUN apk --no-cache add tzdata
WORKDIR /app
COPY --from=build /app/server .
COPY --from=build /app/migrate .
COPY --from=build /app/db ./db
CMD ["./server"]