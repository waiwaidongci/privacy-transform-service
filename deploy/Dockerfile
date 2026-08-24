FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /privacy-transform-service ./cmd/privacy-transform-service
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /privacy-transform-service /privacy-transform-service
EXPOSE 8083
ENTRYPOINT ["/privacy-transform-service"]
