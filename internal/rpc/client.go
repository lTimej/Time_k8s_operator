package rpc

import (
	"Time_k8s_operator/conf"
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	clients = map[string]*grpc.ClientConn{}
	lock    sync.Mutex
)

func GrpcClient(name string) *grpc.ClientConn {
	lock.Lock()
	defer lock.Unlock()
	conn := newGrpcClient()
	if conn == nil {
		log.Fatalln("did not connect")
		return nil
	}
	clients[name] = conn
	return conn
}

func newGrpcClient() *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	conn, err := grpc.DialContext(ctx, conf.GrpcConfig.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil
	}
	return conn
}
