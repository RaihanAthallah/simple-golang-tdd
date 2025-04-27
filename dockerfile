# Gunakan base image Go
FROM golang:1.23.0

# Set working directory
WORKDIR /app

# Copy semua file ke dalam container
COPY . .

# Install dependency
RUN go mod download

# Build binary
RUN go build -o main .

# Expose port (ambil dari .env nanti di compose)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
