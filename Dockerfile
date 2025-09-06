# syntax=docker/dockerfile:1

FROM golang:1.25 AS build

WORKDIR /go/src
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /app

FROM build AS test
RUN go test -v ./...


FROM gcr.io/distroless/base-debian11 AS release
WORKDIR /
COPY --from=build /app /app

USER nonroot:nonroot
ENTRYPOINT ["/app"]
