FROM golang:latest as builder

RUN mkdir /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD /cmd/gen /app/
ADD /pkg /app/pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go

FROM scratch
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]