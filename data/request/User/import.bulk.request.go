package user

type BulkImportUsersRequest struct {
	Users          []ImportedUserData `json:"users" validate:"required,min=1"`
	CustomPassword string             `json:"customPassword"`
}
