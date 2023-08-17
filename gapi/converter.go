package gapi

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	db "github.com/liquiddev99/dropbyte-backend/db/sqlc"
	"github.com/liquiddev99/dropbyte-backend/pb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
