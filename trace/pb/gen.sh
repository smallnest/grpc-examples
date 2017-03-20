GOGO_ROOT=${GOPATH}/src/github.com/gogo/protobuf

protoc -I.:${GOPATH}/src  --gogofaster_out=plugins=grpc:. helloworld.proto
