
# Builder
FROM golang:1.13.4-alpine as builder

RUN apk update && apk upgrade && \
  apk --update add git gcc make

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
  apk --update --no-cache add tzdata && \
  mkdir /app

WORKDIR /app

EXPOSE 9000

COPY --from=builder /app/bid-tracker /app

# Add a user to run nginx in non-root mode
RUN mkdir /user && \
  echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
  echo 'nobody:x:65534:' > /user/group

USER nobody:nobody

CMD /app/bid-tracker
