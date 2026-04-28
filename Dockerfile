# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/uri-template-tester .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/uri-template-tester /uri-template-tester
COPY --from=build /src/public /public/
USER nonroot:nonroot
ENTRYPOINT ["/uri-template-tester"]
EXPOSE 8080
