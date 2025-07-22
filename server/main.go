package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"go-grpc-todo/proto"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type server struct {
	proto.UnimplementedTodoServiceServer
	todo map[string]*proto.Todo
	mu   sync.Mutex
}

func (s *server) AddTodo(ctx context.Context, req *proto.AddTodoRequest) (*proto.AddTodoResponse, error) {
	id := fmt.Sprintf("%d", time.Now().Unix())
	add := proto.Todo{Title: req.Title + "50051", Description: req.Description, Status: proto.Status_TODO_PENDING}
	add.Id = id
	s.mu.Lock()
	if s.todo == nil {
		s.todo = make(map[string]*proto.Todo)
	}
	s.todo[id] = &add
	s.mu.Unlock()
	resp := proto.AddTodoResponse{Todo: &add}
	return &resp, nil
}
func (s *server) ListTodos(ctx context.Context, req *proto.ListTodosRequest) (*proto.ListTodosResponse, error) {
	resp := proto.ListTodosResponse{}
	for _, v := range s.todo {
		resp.Todos = append(resp.Todos, v)
	}
	return &resp, nil
}
func (s *server) CompleteTodo(ctx context.Context, req *proto.CompleteTodoRequest) (*proto.CompleteTodoResponse, error) {
	resp := proto.CompleteTodoResponse{}
	if v, ok := s.todo[req.Id]; ok {
		v.Status = proto.Status_TODO_COMPLETED
		resp.Todo = v
	}
	return &resp, nil
}
func registerService(address string, port int, serviceName string) error {
	config := api.DefaultConfig()
	config.Address = "192.168.5.4:8500" // Consul 容器地址
	client, err := api.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create Consul client: %v", err)
	}

	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, address, port),
		Name:    serviceName,
		Address: address,
		Port:    port,
		Tags:    []string{"grpc", "todo"},
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", address, port, serviceName),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}
	log.Printf("Service %s registered with Consul at %s:%d", serviceName, address, port)
	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	addr := listener.Addr().(*net.TCPAddr)
	host := "192.168.5.4"
	port := addr.Port

	if err := registerService(host, port, "todo-service"); err != nil {
		log.Fatalf("Failed to register with Consul: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterTodoServiceServer(s, &server{})
	// 注册健康检查服务
	healthServer := health.NewServer()
	healthServer.SetServingStatus("todo-service", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	log.Println("Server is running on port :50051")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
