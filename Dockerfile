FROM golang:1.24.2 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /wallet-service cmd/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /wallet-service /
CMD ["/wallet-service"]