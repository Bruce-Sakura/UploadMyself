FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /uploadmyself .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates ffmpeg python3 py3-pip py3-numpy
WORKDIR /app
COPY --from=builder /uploadmyself .
COPY ../ml/scripts /app/ml/scripts
RUN pip3 install --break-system-packages -r /app/ml/scripts/requirements.txt 2>/dev/null || true
RUN mkdir -p /app/uploads
EXPOSE 8000
CMD ["./uploadmyself"]
