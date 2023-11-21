FROM golang

WORKDIR /app/

COPY ./go.mod ./go.sum /app/

RUN go mod download

COPY . /app/

CMD ["go", "run", "./cmd/app"]