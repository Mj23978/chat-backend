package user

import (
	errs "github.com/pkg/errors"
	"gopkg.in/dealancer/validate.v2"

	"github.com/mj23978/chat-backend-x/utils"
	m "github.com/mj23978/chat-backend/pkg/models"
)

type UserService interface {
	SignUp(user *m.User) error
	SignIn(username, password string) (*m.User, error)
	Status(id string, status m.StatusType) error
	FindById(id string) (*m.User, error)
	UpdateUser(user *m.User) error
	Logout()
	Stream()
	Search()
	// Get User Chats Sort By Last Update
	GetChats()
}

type userService struct {
	dbRepo    UserDBRepo
	cacheRepo UserCacheRepo
}

func NewUserService(dbRepo UserDBRepo, cacheRepo UserCacheRepo) UserService {
	return &userService{
		dbRepo:    dbRepo,
		cacheRepo: cacheRepo,
	}
}

func (svc *userService) Logout() {
	panic("implement me")
}

func (svc *userService) Stream() {
	panic("implement me")
}

func (svc *userService) Search() {
	panic("implement me")
}

func (svc *userService) GetChats() {
	panic("implement me")
}

func (svc *userService) SignIn(username, password string) (*m.User, error) {
	res, err := svc.dbRepo.FindBy(map[string]interface{}{
		"username": username,
	})
	if err != nil {
		return nil, errs.Wrap(err, "service.News.GetById")
	}
	if res == nil {
		return nil, errs.New("No User")
	}
	if err := svc.dbRepo.DecryptPassword(res.Password, password); err != nil {
		return nil, errs.Wrap(err, "service.News.GetById")
	}
	return res, nil
}

func (svc *userService) SignUp(user *m.User) error {
	if e := validate.Validate(user); e != nil {
		return errs.Wrap(e, "service.User.SignUp")
	}
	password, err := svc.dbRepo.EncryptPassword(user.Password)
	if err != nil {
		return errs.Wrap(err, "service.User.SignUp")
	}
	user.Password = password
	if err := svc.dbRepo.Create(user); err != nil {
		return errs.Wrap(err, "service.User.SignUp")
	}
	if err := svc.cacheRepo.Create(user); err != nil {
		return errs.Wrap(err, "service.User.SignUp")
	}
	return nil
}

func (svc *userService) Status(id string, status m.StatusType) error {
	if status == m.Online {
		user, err := svc.FindById(id)
		if err != nil {
			return errs.Wrap(err, "service.User.Status")
		}
		user.Status = m.Online
		if err := svc.cacheRepo.Create(user); err != nil {
			return errs.Wrap(err, "service.User.Status")
		}
	} else if status == m.Offline {
		if err := svc.cacheRepo.Delete(id, "username", "status", "level", "token"); err != nil {
			return errs.Wrap(err, "service.User.Status")
		}
	} else if status == m.Playing {
		if err := svc.cacheRepo.Update(id, map[string]interface{}{
			"status": m.Playing,
		}); err != nil {
			return errs.Wrap(err, "service.User.Status")
		}
	}
	return nil
}

func (svc *userService) FindById(id string) (*m.User, error) {
	res, err := svc.dbRepo.FindBy(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, errs.Wrap(err, "service.User.GetById")
	}
	if res == nil {
		return nil, errs.New("No User")
	}
	return res, nil
}

func (svc *userService) UpdateUser(user *m.User) error {
	fields := make(map[string]interface{})
	utils.AppendMapIfExists(utils.StringNullCheck(user.Username), fields, "username", user.Username)
	utils.AppendMapIfExists(utils.StringNullCheck(user.Token), fields, "token", user.Token)
	err := svc.dbRepo.Update(map[string]interface{}{
		"id": user.ID,
	}, fields)
	if err != nil {
		return errs.Wrap(err, "service.User.UpdateUser")
	}
	return nil
}
