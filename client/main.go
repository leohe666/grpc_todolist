package main

import (
	"context"
	"fmt"
	"go-grpc-todo/proto"
	"log"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func discoverService(serviceName string) (string, error) {
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := api.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("failed to create Consul client: %v", err)
	}

	// 查询服务
	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("failed to discover service: %v", err)
	}
	if len(services) == 0 {
		return "", fmt.Errorf("no healthy service instances found for %s", serviceName)
	}

	// 选择第一个健康的服务实例
	service := services[0]
	return fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port), nil
}

func main() {

	// 从 Consul 发现服务地址
	serviceAddr, err := discoverService("todo-service")
	if err != nil {
		log.Fatalf("Service discovery failed: %v", err)
	}
	log.Printf("Discovered service at %s", serviceAddr)

	conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	// respc, err := client.CompleteTodo(ctx, &proto.CompleteTodoRequest{Id: "1753024446"})
	// if err != nil {
	// 	log.Fatal("complete todo error:", err)
	// }
	// fmt.Println("complete todo success response:", respc.Todo)
}
