FROM golang:alpine


ENV GIN_MODE=release
ENV PORT=8080

WORKDIR /go/src/go-docker-dev.to/src

COPY src/ /go/src/go-docker-dev.to/src
#COPY templates /go/src/go-docker-dev.to/templates
RUN go mod download

# Run the two commands below to install git and dependencies for the project. 
RUN apk update && apk add --no-cache git
# RUN go get ./...
RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-redis/redis
RUN go get github.com/nu7hatch/gouuid


RUN go build .

EXPOSE $PORT

ENTRYPOINT ["./shortener"]