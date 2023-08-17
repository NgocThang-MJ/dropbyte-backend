package gapi

import (
	"context"

	"github.com/jackc/pgx/v5"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/liquiddev99/dropbyte-backend/pb"
	"github.com/liquiddev99/dropbyte-backend/util"
	"github.com/liquiddev99/dropbyte-backend/validation"
)

func (server *Server) LoginUser(
	ctx context.Context,
	req *pb.LoginUserRequest,
) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Email doesn't exists: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "Failed to get user: %s", err)
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Password not match: %s", err)
	}

	//	token, err := server.token.CreateToken(user.ID, server.config.AccessTokenDuration)

	rsp := &pb.LoginUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateLoginUserRequest(
	req *pb.LoginUserRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := validation.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
