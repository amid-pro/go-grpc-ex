package main

import (
	"context"
	"encoding/base64"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	port = ":50051"

	isCheckAuth     = true
	loginPassword   = base64.StdEncoding.EncodeToString([]byte("admin:password"))
	errInvalidToken = status.Errorf(codes.Unauthenticated, "Неправильный логин/пароль")

	items = map[uint32]string{
		1:   "Описание товара №1",
		2:   "Описание товара №2",
		444: "Описание товара №444",
	}
)

type server struct {
	UnimplementedItemsServer
}

func checkAuth(ctx *context.Context) bool {
	if !isCheckAuth {
		return true
	}

	md, _ := metadata.FromIncomingContext(*ctx)
	token := strings.TrimPrefix(md["authorization"][0], "Basic ")

	return token == loginPassword
}

// Унарный middleware
func unaryInceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//log.Println(info.FullMethod)

	if !checkAuth(&ctx) {
		return nil, errInvalidToken
	}

	m, err := handler(ctx, req)
	return m, err
}

func (s *server) GetItem(ctx context.Context, itemRequest *ItemId) (*ItemDescription, error) {
	id := itemRequest.GetId()
	value, ok := items[id]

	if ok {
		return &ItemDescription{Description: value}, nil
	}

	return nil, status.Errorf(codes.NotFound, "Не найдено")
}

type ServerStreamInceptor struct {
	grpc.ServerStream
}

func (ssi *ServerStreamInceptor) Rec(m interface{}) error {
	log.Println("Rec: ", m)
	return ssi.ServerStream.RecvMsg(m)
}

func (ssi *ServerStreamInceptor) Send(m interface{}) error {
	log.Println("Send: ", m)
	return ssi.ServerStream.SendMsg(m)
}

func ServerStreamWrapper(s grpc.ServerStream) grpc.ServerStream {
	return &ServerStreamInceptor{s}
}

// Stream middleware
func streamInceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	context := ss.Context()
	if !checkAuth(&context) {
		return errInvalidToken
	}

	err := handler(srv, ServerStreamWrapper(ss))
	return err
}

func (s *server) GetStreamItems(ids *ItemsIds, stream Items_GetStreamItemsServer) error {

	for _, rawId := range ids.GetIds() {

		id := rawId.Id
		value, ok := items[id]

		if ok {
			item := Item{
				Id:          id,
				Description: value,
			}

			time.Sleep(time.Second) // Типа задержка

			stream.Send(&item)
		}
	}

	return nil
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Ошибка слушателя: %v", err)
	}

	gServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInceptor),
		grpc.StreamInterceptor(streamInceptor),
	)

	RegisterItemsServer(gServer, &server{})

	log.Printf("Сервер запущен %v", lis.Addr())

	if err := gServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}

}
