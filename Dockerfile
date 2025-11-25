FROM golang:1.25.4-alpine3.22 AS  base
WORKDIR app

RUN apk update && apk add mesa mesa-dri-gallium

COPY go.* .
RUN go mod download && go mod tidy

COPY . ./
RUN GOOS=linux go build -o covlet

FROM scratch

COPY --from=base covlet covlet

CMD ["covlet"]