package api

import (
	"errors"
	"fmt"
	"grpc_api/database"
	"grpc_api/proto"
	"strconv"

	context "golang.org/x/net/context"
)

type Server struct {
	Database *database.Database
}

func (s *Server) AddUser(ctx context.Context, in *proto.AddUserRequest) (*proto.UserResponse, error) {
	user := &database.User{
		Email:    in.Email,
		Password: in.Password,
		Name:     in.Name,
	}

	resp := &proto.UserResponse{}

	newUser, err := s.Database.AddUser(user)
	if err != nil {
		resp.Success = false
		resp.Data = &proto.UserResponse_Error{Error: &proto.Error{Reason: fmt.Sprintf("%s", err)}}
		return resp, nil
	}

	resp.Success = true
	resp.Data = &proto.UserResponse_User{User: &proto.User{Id: newUser.ID, Email: newUser.Email, Name: newUser.Name}}
	return resp, nil
}

func (s *Server) LoginUser(ctx context.Context, in *proto.LoginUserRequest) (*proto.UserResponse, error) {
	user := &database.User{
		Email:    in.Email,
		Password: in.Password,
	}

	resp := &proto.UserResponse{}

	isValid, err := s.Database.ValidateUser(user)
	if !isValid {
		err = errors.New("INVALID_CREDENTIALS")
	}
	if err != nil {
		resp.Success = false
		resp.Data = &proto.UserResponse_Error{Error: &proto.Error{Reason: fmt.Sprintf("%s", err)}}
		return resp, nil
	}

	dbUser, _ := s.Database.FindUser(user)
	newUser, err := s.Database.CreateToken(dbUser)
	if err != nil {
		resp.Success = false
		resp.Data = &proto.UserResponse_Error{Error: &proto.Error{Reason: fmt.Sprintf("%s", err)}}
		return resp, nil
	}

	resp.Success = true
	resp.Data = &proto.UserResponse_User{User: &proto.User{Id: newUser.ID, Email: newUser.Email, Name: newUser.Name, AccessToken: newUser.AccessToken}}
	return resp, nil
}

func (s *Server) GetUser(ctx context.Context, in *proto.GetUserRequest) (*proto.UserResponse, error) {
	headers := in.GetHeaders()
	if headers["authorization"] == "" {
		return nil, errors.New("NO_AUTH_PROVIDED")
	}

	params := in.GetParams()
	userID := params["id"]
	if userID == "" {
		return nil, errors.New("NO_USER_ID")
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID")
	}
	resp := &proto.UserResponse{}

	user, err := s.Database.FindUser(&database.User{ID: id})
	if err != nil {
		resp.Success = false
		resp.Data = &proto.UserResponse_Error{Error: &proto.Error{Reason: fmt.Sprintf("%s", err)}}
		return resp, nil
	}

	if user.AccessToken != headers["authorization"] {
		err = errors.New("UNAUTHORIZED")
		resp.Success = false
		resp.Data = &proto.UserResponse_Error{Error: &proto.Error{Reason: fmt.Sprintf("%s", err)}}
		return resp, nil
	}

	resp.Success = true
	resp.Data = &proto.UserResponse_User{User: &proto.User{Id: user.ID, Email: user.Email, Name: user.Name}}
	return resp, nil
}
