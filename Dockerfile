FROM golang:1.23.4-alpine as dependencies

WORKDIR /app
COPY go.mod go.sum ./


RUN go mod tidy

COPY . ./
RUN CG0_ENABLE=0 go build -o /main -ldflags="-w -s"


CMD [ "/main" ]