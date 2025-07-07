package tinyurl

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"github.com/evrblk/monstera-example/tinyurl/corepb"
	"github.com/evrblk/monstera-example/tinyurl/gatewaypb"
	monsterax "github.com/evrblk/monstera/x"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TinyUrlServiceApiServer struct {
	gatewaypb.UnimplementedTinyUrlServiceApiServer

	coreApiClient TinyUrlServiceCoreApi
}

func (s *TinyUrlServiceApiServer) Close() {
	log.Println("Stopping ApiServer...")
}

func (s *TinyUrlServiceApiServer) GetUser(ctx context.Context, request *gatewaypb.GetUserRequest) (*gatewaypb.GetUserResponse, error) {
	userId, err := DecodeUserId(request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	resp1, err := s.coreApiClient.GetUser(ctx, &corepb.GetUserRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.GetUserResponse{
		User: userToFront(resp1.User),
	}, nil
}

func (s *TinyUrlServiceApiServer) GetShortUrl(ctx context.Context, request *gatewaypb.GetShortUrlRequest) (*gatewaypb.GetShortUrlResponse, error) {
	shortUrlId, err := DecodeShortUrlId(request.ShortUrl)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid short URL")
	}

	resp1, err := s.coreApiClient.GetShortUrl(ctx, &corepb.GetShortUrlRequest{
		ShortUrlId: shortUrlId,
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.GetShortUrlResponse{
		ShortUrl: shortUrlToFront(resp1.ShortUrl),
	}, nil
}

func (s *TinyUrlServiceApiServer) CreateShortUrl(ctx context.Context, request *gatewaypb.CreateShortUrlRequest) (*gatewaypb.CreateShortUrlResponse, error) {
	userId, err := DecodeUserId(request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	now := time.Now()

	resp1, err := s.coreApiClient.CreateShortUrl(ctx, &corepb.CreateShortUrlRequest{
		ShortUrlId: &corepb.ShortUrlId{
			UserId:     userId,
			ShortUrlId: rand.Uint32(),
		},
		FullUrl: request.FullUrl,
		Now:     now.UnixNano(),
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.CreateShortUrlResponse{
		ShortUrl: shortUrlToFront(resp1.ShortUrl),
	}, nil
}

func NewTinyUrlServiceApiServer(coreApiClient TinyUrlServiceCoreApi) *TinyUrlServiceApiServer {
	return &TinyUrlServiceApiServer{
		coreApiClient: coreApiClient,
	}
}
