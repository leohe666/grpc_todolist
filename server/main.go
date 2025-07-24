package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "go-grpc-todo/proto"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTodoServiceServer
	todo map[string]*pb.Todo
	mu   sync.Mutex
}

func (s *server) AddTodo(ctx context.Context, req *pb.AddTodoRequest) (*pb.AddTodoResponse, error) {
	log.Printf("Received AddTodo request: title=%s, description=%s", req.Title, req.Description)
	id := fmt.Sprintf("%d", time.Now().UnixNano()) // 使用更高精度 ID
	todo := &pb.Todo{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      pb.Status_TODO_PENDING,
	}
	s.mu.Lock()
	if s.todo == nil {
		s.todo = make(map[string]*pb.Todo)
	}
	s.todo[id] = todo
	s.mu.Unlock()
	log.Printf("Added todo: %v", todo)
	return &pb.AddTodoResponse{Todo: todo}, nil
}

func (s *server) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	log.Printf("Received ListTodos request")
	s.mu.Lock()
	todos := make([]*pb.Todo, 0, len(s.todo))
	for _, todo := range s.todo {
		todos = append(todos, todo)
	}
	s.mu.Unlock()
	log.Printf("Listed %d todos", len(todos))
	return &pb.ListTodosResponse{Todos: todos}, nil
}

func (s *server) CompleteTodo(ctx context.Context, req *pb.CompleteTodoRequest) (*pb.CompleteTodoResponse, error) {
	log.Printf("Received CompleteTodo request: id=%s", req.Id)
	s.mu.Lock()
	todo, exists := s.todo[req.Id]
	s.mu.Unlock()
	if !exists {
		log.Printf("Todo %s not found", req.Id)
		return nil, fmt.Errorf("todo %s not found", req.Id)
	}
	todo.Status = pb.Status_TODO_COMPLETED
	log.Printf("Completed todo: %v", todo)
	return &pb.CompleteTodoResponse{Todo: todo}, nil
}

func initEtcdClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
}

func registerService(cli *clientv3.Client, serviceName, addr string) error {
	mg, err := endpoints.NewManager(cli, serviceName)
	if err != nil {
		return fmt.Errorf("failed to create endpoints manager: %v", err)
	}
	leaseResp, err := cli.Grant(context.Background(), 10)
	if err != nil {
		return fmt.Errorf("failed to grant lease: %v", err)
	}
	err = mg.AddEndpoint(context.Background(), serviceName+"/"+addr, endpoints.Endpoint{Addr: addr}, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("failed to add endpoint: %v", err)
	}
	keepAliveCh, err := cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return fmt.Errorf("failed to keep alive: %v", err)
	}
	go func() {
		for range keepAliveCh {
			// log.Printf("Lease %d kept alive for %s", leaseResp.ID, addr)
		}
	}()
	return nil
}

func main() {
	port := flag.String("port", "50052", "gRPC server port")
	flag.Parse()

	addr := fmt.Sprintf("127.0.0.1:%s", *port)
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &server{})

	cli, err := initEtcdClient()
	if err != nil {
		log.Fatalf("failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	if err := registerService(cli, "nuts/v1/EchoService", addr); err != nil {
		log.Fatalf("failed to register service: %v", err)
	}

	log.Printf("Server is running on %s", addr)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
