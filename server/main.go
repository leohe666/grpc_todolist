package main

import (
	"context"
	"fmt"
	"go-grpc-todo/proto"
	"log"
	"net"
	"sync"
	"time"

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

	log.Println("Server is running on port :50052")

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
