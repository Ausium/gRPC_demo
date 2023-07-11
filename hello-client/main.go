package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc_demo/hello-server/proto"
	"time"
)

// 自定义token认证代码
type ClientTokenAuth struct {
}

func (c ClientTokenAuth) GetRequestMetadata(ctx context.Context, url ...string) (map[string]string, error) {
	return map[string]string{
		"appId":  "grpcstudy",
		"appKey": "123456",
	}, nil
}

func (c ClientTokenAuth) RequireTransportSecurity() bool {
	return false
}

func main() {

	start := time.Now() // 记录开始时间

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithPerRPCCredentials(new(ClientTokenAuth)))

	//creds, _ := credentials.NewClientTLSFromFile("D:\\Environment\\GoWorks\\gRPC_demo\\key\\test.pem", "*.grpcstudy.com")

	//连接到server端，此处禁用安全传输，没有加密和验证`
	conn, err := grpc.Dial("127.0.0.1:9090", opts...)
	//conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}

	defer conn.Close()
	//建立连接
	client := pb.NewSayHelloClient(conn)
	//执行grpc调用（这个方法在服务端来实现并返回结果）
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "chen"})

	fmt.Printf("返回的结果为：%v\n", resp)
	endTime := time.Since(start)
	fmt.Printf("程序耗时时间为: %v", endTime)
}
