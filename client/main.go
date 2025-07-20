package main

import (
	"context"
	"fmt"
	"go-grpc-todo/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
