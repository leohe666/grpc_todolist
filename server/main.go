package main

import (
	"context"
	"fmt"
	"go-grpc-todo/proto"
	"log"
	"net"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedTodoServiceServer
	todo map[string]*proto.Todo
	mu   sync.Mutex
}

func (s *server) AddTodo(ctx context.Context, req *proto.AddTodoRequest) (*proto.AddTodoResponse, error) {
	id := fmt.Sprintf("%d", time.Now().Unix())
	add := proto.Todo{Title: req.Title, Description: req.Description, Status: proto.Status_TODO_PENDING}
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

func initEtcdClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // etcd 服务地址
		DialTimeout: 5 * time.Second,
	})
}
func registerService(cli *clientv3.Client, serviceName, addr string, ttl int64) error {
	// 创建租约
	leaseResp, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	// 注册服务地址
	key := "/services/" + serviceName + "/" + addr
	_, err = cli.Put(context.Background(), key, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	// 保持租约活跃（心跳）
	keepAliveCh, err := cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return err
	}

	// 异步处理心跳响应
	go func() {
		for range keepAliveCh {
			// 持续保持租约活跃
		}
	}()

	return nil
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
func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterTodoServiceServer(s, &server{})

	// / 初始化 etcd 客户端
	cli, err := initEtcdClient()
	if err != nil {
		log.Fatalf("failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	// 注册服务
	serviceName := "todolist"
	addr := "localhost:50052"
	if err := registerService(cli, serviceName, addr, 10); err != nil {
		log.Fatalf("failed to register service: %v", err)
	}

	log.Println("Server is running on port :50052")

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
