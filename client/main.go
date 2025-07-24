package main

import (
	"context"
	"fmt"
	"go-grpc-todo/proto"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initEtcdClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // etcd 服务地址
		DialTimeout: 5 * time.Second,
	})
}
func discoverService(cli *clientv3.Client, serviceName string) ([]string, error) {
	prefix := "/services/" + serviceName + "/"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, kv := range resp.Kvs {
		addrs = append(addrs, string(kv.Value))
	}
	return addrs, nil
}

func main() {
	// 初始化 etcd 客户端
	cli, err := initEtcdClient()
	if err != nil {
		log.Fatalf("failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	// 发现服务
	addrs, err := discoverService(cli, "todolist")
	if err != nil {
		log.Fatalf("failed to discover service: %v", err)
	}

	conn, err := grpc.NewClient(addrs[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("client error:", err)
	}
	defer conn.Close()

	client := proto.NewTodoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.AddTodo(ctx, &proto.AddTodoRequest{Title: "todo 2", Description: "todo 2"})
	if err != nil {
		log.Fatal("add todo error:", err)
	}
	fmt.Println("add todo success response:", resp.Todo)

	respl, err := client.ListTodos(ctx, &proto.ListTodosRequest{})
	if err != nil {
		log.Fatal("list todo error:", err)
	}
	fmt.Println("list todo success response:", respl.Todos)

	respc, err := client.CompleteTodo(ctx, &proto.CompleteTodoRequest{Id: "1753024446"})
	if err != nil {
		log.Fatal("complete todo error:", err)
	}
	fmt.Println("complete todo success response:", respc.Todo)
}
