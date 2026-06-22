package user

import (
	request "Microservice/data/request/User"
	response "Microservice/data/response/User"
	"Microservice/helper"
	"mime/multipart"
)

type UserService interface {
	Create(data request.CreateUserRequest, orgID string) *helper.ErrorModel
	Get(id string, orgID string) (*response.UserResponse, *helper.ErrorModel)
	GetAll(orgID string) ([]response.UserResponse, *helper.ErrorModel)
	GetAllUserExceptCurrent(userId string, orgID string) ([]response.UserResponse, *helper.ErrorModel)
	Update(request request.UpdateUserRequest, orgID string) *helper.ErrorModel
	Delete(id string, orgID string) *helper.ErrorModel
	MultipleDelete(ids []string, orgID string) *helper.ErrorModel

	UpdateBiodata(id string, request request.UpdateBiodataRequest, orgID string) *helper.ErrorModel
	UpdateEmail(id string, request request.UpdateEmailRequest, orgID string) *helper.ErrorModel
	UpdateRole(request request.UpdateRoleRequest, orgID string) *helper.ErrorModel
	UpdatePassword(request request.UpdatePasswordRequest, orgID string) *helper.ErrorModel
	UpdateAccess(request request.UpdateAccessRequest, orgID string) *helper.ErrorModel

	PreviewImport(file *multipart.FileHeader, columnMappingJSON string, orgID string) (*response.PreviewImportResponse, *helper.ErrorModel)
	BulkImport(request request.BulkImportUsersRequest, orgID string) (*response.BulkImportResponse, *helper.ErrorModel)
	UnlockUser(userId string, orgID string) *helper.ErrorModel
}
