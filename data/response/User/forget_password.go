package response

type VerifyForgetPassword struct {
	Registered bool `json:"registered"`
}

type ResetPassword struct {
	PasswordValid bool `json:"registered"`
}
