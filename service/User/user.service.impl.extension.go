package user

import (
	response "Microservice/data/response/User"
	"Microservice/model"
)

func (t UserServiceImpl) mapUsertoUserResponse(users []model.User) []response.UserResponse {
	responseUser := make([]response.UserResponse, len(users))
	for i, user := range users {
		responseUser[i] = t.convertUserToUserResponse(user)
	}
	return responseUser
}

func (t UserServiceImpl) convertUserToUserResponse(user model.User) response.UserResponse {
	// Perform necessary conversion logic here, potentially selecting specific fields
	response := response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Position:  user.Position,
		Access:    user.Access,
		Phone:     user.Phone,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}

	return response
}
