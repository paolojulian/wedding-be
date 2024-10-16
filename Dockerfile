FROM golang:1.22-bookworm

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
# Set environment variables to ensure the binary is built for Linux amd64 (which is what Google Cloud Run expects)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /usr/local/bin/app

EXPOSE 8080

CMD ["/usr/local/bin/app"]