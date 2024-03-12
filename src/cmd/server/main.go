package main

import (
	"context"
	"log"
	"net"

	user_managementpb "github.com/dkpcb/user-management-service/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// サーバーの実装
type server struct {
	user_managementpb.UnimplementedUserManagementServiceServer
	users map[string]*user_managementpb.User
}

// ユーザー追加
func (s *server) AddUser(ctx context.Context, in *user_managementpb.AddUserRequest) (*user_managementpb.AddUserResponse, error) {
	id := in.GetUser().GetId() // 一意なIDを想定
	s.users[id] = in.GetUser()
	return &user_managementpb.AddUserResponse{Id: id}, nil
}

// ユーザー取得
func (s *server) GetUser(ctx context.Context, in *user_managementpb.GetUserRequest) (*user_managementpb.GetUserResponse, error) {
	user, exists := s.users[in.GetId()]
	if !exists {
		// gRPCエラーとしてNotFoundエラーを返す
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	return &user_managementpb.GetUserResponse{User: user}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// ここで、サービスをサーバーに登録します。
	user_managementpb.RegisterUserManagementServiceServer(s, &server{users: make(map[string]*user_managementpb.User)})

	// サーバーにリフレクションサービスを登録します。
	reflection.Register(s)

	// サーバーを起動します。
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
