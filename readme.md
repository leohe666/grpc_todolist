protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto

go grpc consul

docker pull hashicorp/consul:1.21
 
docker run --name consul-server -d -p 8500:8500 -p 8600:8600/udp hashicorp/consul agent -server -bootstrap-expect 1 -ui -client=0.0.0.0

打开浏览器访问 http://localhost:8500