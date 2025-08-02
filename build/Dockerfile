FROM mcr.microsoft.com/devcontainers/go:1.24-bookworm
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-go-application

CMD ["/docker-go-application"]