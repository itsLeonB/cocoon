FROM golang:1.24-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -buildvcs=false -ldflags='-w -s' \
    -o /cocoon ./cmd/grpc/main.go

FROM gcr.io/distroless/static-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /cocoon /cocoon

EXPOSE 50051

USER nonroot:nonroot

ENTRYPOINT ["/cocoon"]
