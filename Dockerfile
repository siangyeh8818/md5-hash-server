FROM golang:1.16.4-stretch as builder

COPY go.mod /go/src/github.com/siangyeh8818/md5-hash-server/go.mod
COPY go.sum /go/src/github.com/siangyeh8818/md5-hash-server/go.sum

# Run golang at any directory, not neccessary $GOROOT, $GOPATH
ENV GO111MODULE=on
WORKDIR /go/src/github.com/siangyeh8818/md5-hash-server

# RUN go mod init github.com/pnetwork/sre.monitor.metrics
RUN go mod download
COPY main.go /go/src/github.com/siangyeh8818/md5-hash-server/
COPY internal /go/src/github.com/siangyeh8818/md5-hash-server/internal
#COPY pkg /go/src/github.com/pnetwork/sre.monitor.metrics/pkg

# Build the Go app
RUN env GOOS=linux GOARCH=amd64 go build -o md5-hash-server -v -ldflags "-s" github.com/siangyeh8818/md5-hash-server/main.go

##### To reduce the final image size, start a new stage with alpine from scratch #####
FROM centos:7

# Run as root
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/siangyeh8818/md5-hash-server/md5-hash-server /usr/local/bin/md5-hash-server

# EXPOSE 8081

ENTRYPOINT [ "md5-hash-server" ] 