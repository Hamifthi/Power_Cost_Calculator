FROM golang:1.17

RUN apt-get update

RUN apt-get install unzip -y

RUN groupadd --gid 1000 data \
  && useradd --uid 1000 --gid data data

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN mkdir inputs, generated_data, outputs
RUN chown -R data:data /usr/src/app/outputs/

ADD inputs.zip ./inputs.zip
RUN unzip ./inputs.zip -d ./inputs

COPY sample.env .env
COPY . .

RUN go build -v -o ./app ./cmd/main.go
CMD go test -v ./...; ./app