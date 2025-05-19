package usecase

import (
	"errors"
	"gatherly-app/models/dto"
	"gatherly-app/repositories"
	"gatherly-app/service"
	"gatherly-app/utils"

	"github.com/gin-gonic/gin"
)

type AuthenticationUseCase interface {
	Login(ctx *gin.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
}

type authenticationUseCase struct {
	userRepo   repositories.UserRepository
	jwtService service.JwtService
}

func NewAuthenticationUseCase(userRepo repositories.UserRepository, jwtService service.JwtService) AuthenticationUseCase {
	return &authenticationUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (uc *authenticationUseCase) Login(ctx *gin.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := uc.userRepo.FindByEmail(request.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := user.CheckPassword(request.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Dapatkan IP address dari context
	ipAddress := ctx.ClientIP()

	// Dapatkan koordinat dari IP
	var geo *utils.Geocode
	geo, err = utils.GetCoordinateFromIP(ipAddress)
	if err != nil {
		// Jika gagal, lanjutkan tanpa koordinat
		geo = &utils.Geocode{}
	}

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email, user.Role, geo.Latitude, geo.Longitude)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Age:      user.Age,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
		Latitude:  geo.Latitude,
		Longitude: geo.Longitude,
	}, nil
}
