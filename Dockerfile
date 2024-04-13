FROM golang:1.21.6-alpine AS BUILDER
WORKDIR /biggest-change
COPY . .                                    
RUN CGO_ENABLED=0 GOOS=linux go build -o biggest-change ./cmd/app/main.go

FROM alpine:latest
WORKDIR /biggest-change
COPY --from=BUILDER /biggest-change ./
EXPOSE 8080

CMD ["./biggest-change"]