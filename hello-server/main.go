package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	pb "grpc_demo/hello-server/proto"
	"net"
)

type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {

	//获取元数据的信息
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("未传输token！")
	}
	var appId string
	var appKey string
	if v, ok := md["appid"]; ok {
		appId = v[0]
		fmt.Println("appId", appId)
	}

	if v, ok := md["appkey"]; ok {
		appKey = v[0]
		fmt.Println("appKey", appKey)
	}

	//正常这里也是从数据库中查找，每个用户都会有一个自己的appId
	//然后根据用户的id来查找 appId，然后再校验这appId是否正确，不正确就不让通过
	if appId != "grpcstudy" || appKey != "123456" {
		return nil, errors.New("token校验失败！")
	}

	fmt.Println("hello" + req.RequestName)
	return &pb.HelloResponse{ResponseMsg: "hello" + req.RequestName}, nil
}

func main() {

	//TSL认证
	//两个参数分别是 cretFile, keyFile
	//自签名证书文件和私钥文件
	//creds, _ := credentials.NewServerTLSFromFile("D:\\Environment\\GoWorks\\gRPC_demo\\key\\test.pem", "D:\\Environment\\GoWorks\\gRPC_demo\\key\\test.key")

	//开启端口
	listen, _ := net.Listen("tcp", ":9090")

	//创建grpc服务
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	//在grpc服务端中去注册我们自己编写的服务
	pb.RegisterSayHelloServer(grpcServer, &server{})

	//启动服务
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("failed to server %v", err)
		return
	}
}
