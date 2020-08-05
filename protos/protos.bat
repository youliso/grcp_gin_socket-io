# 编译google.api
protoc -I=./protos --go_out=./protos ./protos/google/api/*.proto

# 编译hello.proto gateway
protoc -I=./protos --grpc-gateway_out=logtostderr=true:./protos ./protos/hello/*.proto

# 编译hello.proto
protoc -I=./protos --go_out=Mgoogle/api/annotations.proto=grpc/protos/google/api/annotations,plugins=grpc:./protos ./protos/hello/*.proto