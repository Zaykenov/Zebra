package service

import (
	"errors"
	"time"
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt             = "hjqrhjqw124617ajfhajs"
	signingKey       = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL         = 21 * time.Hour
	tokenYear        = 365 * 24 * time.Hour
	liitleTokenTTL   = 30 * time.Second
	boymanLocalToken = "ba2228a00d21e19c23e4f210a5b8a300"
)

type tokenClaims struct {
	jwt.StandardClaims
	Role     string `json:"role"`
	UserID   int    `json:"user_id"`
	Shops    []int  `json:"shops"`
	BindShop int    `json:"bind_shop"`
}

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) AddUser(user *model.UserRequest) (int, error) {
	newUser := model.User{
		DeviceID: user.DeviceID,
		Phone:    user.Phone,
		Name:     user.Name,
		Balance:  0,
	}
	if newUser.Phone == "" {
		newUser.Discount = utils.DefaultDiscount
	} else {
		newUser.Discount = utils.TelephoneNumberDiscount
	}
	return s.repo.AddUser(&newUser)
}

func (s *UserService) GetUser(id int) (*model.User, error) {
	return s.repo.GetUser(id)
}

func (s *UserService) GetCurrentOrders(id int) ([]*model.Check, error) {
	return s.repo.GetCurrentOrders(id)
}

func (s *UserService) AddFeedback(feedback *model.FeedbackRequest) error {
	newFeedback := model.Feedback{
		CheckID:      feedback.CheckID,
		ScoreQuality: feedback.ScoreQuality,
		ScoreService: feedback.ScoreService,
		Feedback:     feedback.Feedback,
	}
	return s.repo.AddFeedback(&newFeedback)
}

func (s *UserService) GenCode(id int) (string, int64, error) {
	userCode, err := s.repo.GetUserCode(id)
	if err != nil {
		if err.Error() != "record not found" {
			return "", -1, err
		}
	}
	if userCode != nil {
		if userCode.ExpireTime >= time.Now().Local().Unix() {
			return userCode.Code, userCode.ExpireTime, nil
		}
	}
	var code string
	for {
		code = utils.GenCode()
		_, err = s.repo.CheckCode(code)
		if err != nil {
			if err.Error() == "code exists" {
				continue
			}
			if err.Error() != "record not found" {
				return "", -1, err
			}
		}

		userCode, err = s.repo.SetCode(id, code)
		if err != nil {
			return "", -1, err
		}
		break
	}
	return code, userCode.ExpireTime, nil
}

func (s *UserService) GetUserByCode(code string) (*model.User, error) {
	return s.repo.GetUserByCode(code)
}

func (s *UserService) GetWorkerByUsername(username string) (*model.Worker, error) {

	res, err := s.repo.GetWorkerByUsername(username)
	if err != nil {
		return nil, err
	}
	res.Token, err = s.GenerateToken(res.ID, res.Role, res.Shops, res.BindShop)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *UserService) CreateWorker(worker *model.ReqWorkerRegistration) error {
	newWorker := model.Worker{
		Name:     worker.Name,
		Phone:    worker.Phone,
		Username: worker.Username,
		Password: worker.Password,
		Role:     worker.Role,
		Shops:    worker.Shops,
		BindShop: worker.Shops[0],
	}
	res, err := s.repo.CreateWorker(&newWorker)

	if err != nil {
		return err
	}
	token, err := s.GenerateToken(res.ID, newWorker.Role, worker.Shops, worker.BindShop)

	if err != nil {
		return err
	}

	res.Token = token

	return s.repo.UpdateWorker(res)
}

func (s *UserService) GenerateToken(userID int, role string, shops []int, bindShop int) (string, error) {
	var expireTime int64
	if role == "worker" {
		expireTime = time.Now().Local().Add(tokenTTL).Unix()
	} else {
		expireTime = time.Now().Local().Add(tokenYear).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Local().Unix(),
			ExpiresAt: expireTime,
		},
		role,
		userID,
		shops,
		bindShop,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *UserService) ParseToken(accessToken string) (int, string, []int, int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", nil, 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", nil, 0, errors.New("token claims are not of type *tokenClaims")
	}
	// if claims.ExpiresAt < time.Now().Local().Unix() || claims.IssuedAt < time.Now().Local().Add(-tokenTTL).Unix() {
	// 	return 0, "", nil, 0, errors.New("token is expired")
	// }

	return claims.UserID, claims.Role, claims.Shops, claims.BindShop, nil
}
