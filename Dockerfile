FROM gocv/opencv:latest as build

WORKDIR /source

ARG GOVERSION=1.22.0

RUN wget https://go.dev/dl/go$GOVERSION.linux-amd64.tar.gz
RUN tar -C /opt -xzf go$GOVERSION.linux-amd64.tar.gz
RUN rm -rf go$GOVERSION.linux-amd64.tar.gz 
ENV PATH="/opt/go/bin:${PATH}"

COPY ./go.mod ./go.sum ./
RUN go mod download && mkdir build

COPY . .

RUN go build -o build ./...

FROM gocv/opencv:latest

WORKDIR /app

COPY --from=build /source/build/camera /app

ENTRYPOINT ["/app/camera"]
