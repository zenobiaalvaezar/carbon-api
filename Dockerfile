FROM golang:1.22 AS builder

# Set the working directory in the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Set environment variables
ENV API_PORT=8080
ENV PORT=8080
ENV DATABASE_URL=postgresql://postgres.vqgvmkrcnphmmnclvwjs:AO7dSLJEBdiBJIDE@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres
ENV REDIS_URL=redis://default:PRVTGXxETtdNYEEVpMApzGzxAVbfQqqm@monorail.proxy.rlwy.net:18146
ENV SMTP_HOST=smtp.gmail.com
ENV SMTP_PORT=587
ENV SMTP_USERNAME=rehanadipurwana@gmail.com
ENV SMTP_PASSWORD=vcjwrthkeoivbpdf
ENV XENDIT_API_URL="https://api.xendit.co/v2/invoices"
ENV XENDIT_SECRET_KEY="xnd_development_T0trkUor1Wf5EdSyEdaPEbTesSKYCt4tDHyX38wddSQ04TlRbvp00UK1i0v5Ql3Q"
ENV MONGODB_URI="mongodb+srv://fr081938:12345678___@cluster0.7d2ce.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
ENV MONGO_DATABASE=carbon
ENV MONGO_COLLECTION=payment_methods

# Build the Go application
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]
