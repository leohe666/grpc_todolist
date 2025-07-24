package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "go-grpc-todo/proto"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

func initEtcdClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
}

func discoverService(cli *clientv3.Client, serviceName string) ([]string, error) {
	prefix := serviceName + "/"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get service addresses: %v", err)
	}

	var addrs []string
	for _, kv := range resp.Kvs {
		addrs = append(addrs, string(kv.Value))
	}
	return addrs, nil
}

func main() {
	// 启用 gRPC 日志
	// grpclog.SetLoggerV2(grpclog.NewLoggerV2(log.Writer(), log.Writer(), log.Writer()))

	// 初始化 etcd 客户端
	cli, err := initEtcdClient()
	if err != nil {
		log.Fatalf("failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	// 调试：手动发现服务地址
	addrs, err := discoverService(cli, "nuts/v1/EchoService")
	if err != nil {
		log.Fatalf("failed to discover service: %v", err)
	}
	if len(addrs) == 0 {
		log.Fatal("no service addresses found")
	}
	log.Printf("Discovered addresses: %v", addrs)

	// 创建 etcd resolver
	etcdResolver, err := resolver.NewBuilder(cli)
	if err != nil {
		log.Fatalf("failed to create resolver: %v", err)
	}

	// 增加拨号超时
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer dialCancel()

	// 使用 etcd resolver 连接 gRPC 服务
	conn, err := grpc.DialContext(
		dialCtx,
		"etcd:///nuts/v1/EchoService",
		grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithMax(3),
			grpc_retry.WithCodes(codes.Unavailable, codes.DeadlineExceeded),
		)),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// 强制触发连接
	// conn.Connect()
	// // 等待连接就绪
	// for i := 0; i < 10; i++ {
	// 	state := conn.GetState()
	// 	log.Printf("gRPC connection state: %v", state)
	// 	if state == connectivity.Ready {
	// 		break
	// 	}
	// 	if !conn.WaitForStateChange(dialCtx, state) {
	// 		log.Printf("connection timed out, final state: %v", state)
	// 		return
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }

	client := pb.NewTodoServiceClient(conn)

	// 测试 AddTodo
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := client.AddTodo(ctx, &pb.AddTodoRequest{Title: "todo 2", Description: "todo 2"})
	if err != nil {
		log.Printf("add todo error: %v", err)
	} else {
		log.Printf("add todo success response: %v", resp.Todo)
	}

	// 测试 ListTodos
	respl, err := client.ListTodos(ctx, &pb.ListTodosRequest{})
	if err != nil {
		log.Printf("list todo error: %v", err)
	} else {
		log.Printf("list todo success response: %v", respl.Todos)
	}

	// // 测试 CompleteTodo
	// if resp != nil && resp.Todo != nil {
	// 	respc, err := client.CompleteTodo(ctx, &pb.CompleteTodoRequest{Id: resp.Todo.Id})
	// 	if err != nil {
	// 		log.Printf("complete todo error: %v", err)
	// 	} else {
	// 		log.Printf("complete todo success response: %v", respc.Todo)
	// 	}
	// } else {
	// 	log.Println("Skipping CompleteTodo: no valid todo ID")
	// }
}
