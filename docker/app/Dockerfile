FROM golang:1.24

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy all the files from the root directory to the /app directory in the container
COPY . ./

RUN go build -o main ./cmd/main.go

EXPOSE 1000

CMD ["./main"]
