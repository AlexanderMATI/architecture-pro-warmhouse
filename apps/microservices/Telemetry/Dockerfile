FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY app.py .

RUN pip install flask

EXPOSE 5002

CMD ["./telemetry.go"]
