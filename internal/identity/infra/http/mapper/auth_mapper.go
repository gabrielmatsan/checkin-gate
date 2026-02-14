package mapper

import (
	authenticatewithgoogle "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/authenticate_with_google"
	getgoogleauthurl "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/get_google_auth_url"
	refreshtoken "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/refresh_token"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/http/dto"
)

func GoogleCallbackRequestToInput(req *dto.GoogleCallbackRequest, ipAddress, userAgent string) *authenticatewithgoogle.Input {
	return &authenticatewithgoogle.Input{
		Code:      req.Code,
		IpAddress: ipAddress,
		UserAgent: userAgent,
	}
}

func AuthOutputToResponse(output *authenticatewithgoogle.Output) *dto.GoogleCallbackResponse {
	return &dto.GoogleCallbackResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		User:         UserToResponse(output.User),
	}
}

func UserToResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.FirstName + " " + user.LastName,
		Role:  string(user.Role),
	}
}

func RefreshTokenRequestToInput(req *dto.RefreshTokenRequest, ipAddress, userAgent string) *refreshtoken.Input {
	return &refreshtoken.Input{
		RefreshToken: req.RefreshToken,
		IpAddress:    ipAddress,
		UserAgent:    userAgent,
	}
}

func RefreshTokenOutputToResponse(output *refreshtoken.Output) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}
}

func GoogleAuthURLOutputToResponse(output *getgoogleauthurl.Output) *dto.GoogleAuthURLResponse {
	return &dto.GoogleAuthURLResponse{
		URL:   output.URL,
		State: output.State,
	}
}
