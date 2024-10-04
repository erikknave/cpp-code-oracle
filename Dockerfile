# Stage 1: Build the Go app
FROM golang:1.22 AS builder
WORKDIR /app
RUN go install github.com/a-h/templ/cmd/templ@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix 'static' -o main .

# Stage 2: Create the final image
# FROM scratch
FROM gcr.io/distroless/static-debian11:latest
WORKDIR /app_prod
COPY --from=builder /app/main .

# Command to run the executable, using wait-for-it to wait for db service
CMD ["/app_prod/main"]
