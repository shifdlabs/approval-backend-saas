package controller

import (
	request "Microservice/data/request/Authentication"
	userResponse "Microservice/data/response/User"
	"Microservice/helper"
	service "Microservice/service/Token"
	"Microservice/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
	tokenService service.TokenService
}

func NewTokenController(tokenService service.TokenService) *TokenController {
	return &TokenController{tokenService: tokenService}
}

func (controller *TokenController) RefreshAccessToken(ctx *gin.Context) {
	var accessToken string
	access := ctx.GetHeader("Authorization")
	if strings.HasPrefix(access, "Bearer ") {
		accessToken = strings.TrimPrefix(access, "Bearer ")
	}
	helper.PrintValue(access, "Access Token")

	// refresh, err := ctx.Cookie("refreshToken")
	// helper.PrintValue(refresh, "Refresh Token")
	// for _, cookie := range ctx.Request.Cookies() {
	// 	fmt.Println("Cookie:", cookie.Name, cookie.Value)
	// }

	var payload request.RefreshAccessTokenRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	helper.PrintValue(payload.RefreshToken, "Payload Refresh Token")

	newAccessToken, newRefreshToken, err := helper.ValidateOrRefreshAccess(accessToken, payload.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	helper.PrintValue(newAccessToken, "New Access Token")
	if newRefreshToken != "" {
		// rotate cookie if new refresh issued
		ctx.SetCookie("refreshToken", newRefreshToken, int(helper.RefreshTTL.Seconds()), "/", "", true, true)
	}

	// userId, errTokenClaims := config.RedisClient.Get(ctx, *identifier).Result()

	// if errTokenClaims != nil {
	// 	msg := "Refresh token not found"
	// 	utils.ErrorResponse(ctx, helper.ErrorModel{Code: 403, Message: msg})
	// }

	// contextID := context.TODO()
	// now := time.Now()

	// errAccess := config.RedisClient.Set(contextID, result.AccessToken.Identifier, result.AccessToken.UserID, time.Unix(*result.AccessToken.ExpiresIn, 0).Sub(now)).Err()
	// if errAccess != nil {
	// 	msg := "Bad Request"
	// 	utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	// }

	// ctx.SetCookie("access_token", *result.AccessToken.Token, 3600, "/", "localhost", false, true)

	utils.SuccessResponse(ctx, userResponse.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}
