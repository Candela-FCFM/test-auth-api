FROM golang:1.24 AS build

WORKDIR /app

ARG ARG_REPO_PACK 
ARG VERSION

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOSUMDB=off \
    GOPROXY="http://${ARG_REPO_PACK}:8081/repository/go-main-proxy/,direct" 

COPY go.* .

RUN go mod tidy && go mod download

COPY . .

RUN go build -o main-api -v -x -trimpath \
  -ldflags="-s -w \
  -X main.version=0.1.0" 

FROM gcr.io/distroless/static-debian12:nonroot AS main-api

WORKDIR /app

COPY --from=build /app/main-api .

CMD ["main-api"]