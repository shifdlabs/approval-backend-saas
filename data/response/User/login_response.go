package response

type LoginResponse struct {
	AccessToken      string    `json:"accessToken"`
	RefreshToken     string    `json:"refreshToken"`
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	Access           bool      `json:"access"`
	UserAbilityRules []Ability `json:"userAbilityRules"`
	Role             int       `json:"role"`
	JobPosition      string    `json:"jobPosition"`
}

type Ability struct {
	Action  string `json:"action"`
	Subject string `json:"subject"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
