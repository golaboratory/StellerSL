FROM golang:1.24-alpine AS dev

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

# For hot-reloading (optional but recommended for development)
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
