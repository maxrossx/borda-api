package service

import (
	"borda/internal/config"
	"borda/internal/domain"
	"borda/internal/repository"
	"borda/pkg/hash"
	"fmt"

	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
	hasher   hash.PasswordHasher
}

func NewAuthService(ur repository.UserRepository, tr repository.TeamRepository,
	hasher hash.PasswordHasher) *AuthService {

	return &AuthService{
		userRepo: ur,
		teamRepo: tr,
		hasher:   hasher,
	}
}

func (s *AuthService) SignUp(input domain.UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	// TODO:
	//		Attach user to the team.
	//		If parsing token or creating new team fails, user should't be created.
	// 		To achive prosess  should be run in transaction.
	userId, err := s.userRepo.SaveUser(input.Username, passwordHash, input.Contact)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return err
		}
		return err
	}

	switch input.AttachTeamMethod {
	case "create":
		if _, err := s.teamRepo.SaveTeam(userId, input.AttachTeamAttribute); err != nil {
			return err
		}
	case "join":
		team, err := s.teamRepo.GetTeamByToken(input.AttachTeamAttribute)
		if err != nil {
			return err
		}

		if err := s.teamRepo.AddMember(team.Id, userId); err != nil {
			return err
		}
	}

	return nil
}

func (s *AuthService) SignIn(input domain.UserSignInInput) (string, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	fmt.Println(passwordHash)

	user, err := s.userRepo.GetUserByCredentials(input.Username, passwordHash)
	if err != nil {
		return "", err
	}

	jwtConf := config.JWT()

	fmt.Println(jwtConf)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(jwtConf.ExpireTime).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.Itoa(user.Id),
	})

	return token.SignedString([]byte(jwtConf.SigningKey))
}
