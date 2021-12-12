package cassandra

//
//import (
//	"context"
//	"github.com/pkg/errors"
//	"golang.org/x/crypto/bcrypt"
//	"time"
//
//	log "github.com/mj23978/chat-backend-x/logger/zerolog"
//	"github.com/mj23978/chat-backend/pkg/domain/user"
//	m "github.com/mj23978/chat-backend/pkg/models"
//)
//
//type userDBRepo struct {
//	client  *Cassandra
//	timeout time.Duration
//}
//
//func UserCassandraRepository(cas *Cassandra) (user.UserDBRepo, error) {
//	repo := &userDBRepo{
//		client:  cas,
//		timeout: time.Second * 10,
//	}
//	return repo, nil
//}
//
//func (r *userDBRepo) Create(user *m.User) error {
//	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
//	defer cancel()
//	//err := ses.Query("INSERT ").Exec()
//	//if err != nil {
//	//	return nil, utils.ErrorWrap(err, "Ko3sger")
//	//}
//	collection := r.client.Database(r.database).Collection("users")
//	res, err := collection.InsertOne(
//		ctx,
//		bson.M{
//			"id":       user.ID,
//			"username": user.Username,
//			"password": user.Password,
//			"email":    user.Email,
//			"token":    user.Token,
//		},
//	)
//	if err != nil {
//		return errors.Wrap(err, "repository.User.Start")
//	}
//	log.Infof("Mongo Create %v", res)
//	return nil
//}
//
//func (r *userDBRepo) FindBy(filter map[string]interface{}) (*m.User, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
//	defer cancel()
//	collection := r.client.Database(r.database).Collection("users")
//	res := collection.FindOne(
//		ctx,
//		filter,
//	)
//	if err := res.Err(); err != nil {
//		return nil, errors.Wrap(err, "repository.User.FindBy")
//	}
//	user := &m.User{}
//	bs, err := res.DecodeBytes()
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return nil, nil
//		}
//		return nil, errors.Wrap(err, "repository.User.FindBy")
//	}
//	err2 := bson.Unmarshal(bs, user)
//	if err2 != nil {
//		return nil, errors.Wrap(err2, "repository.User.FindBy")
//	}
//	return user, nil
//}
//
//func (r *userDBRepo) Delete(filter map[string]interface{}) error {
//	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
//	defer cancel()
//	collection := r.client.Database(r.database).Collection("users")
//	res, err := collection.DeleteOne(
//		ctx,
//		filter,
//	)
//	if err != nil {
//		return errors.Wrap(err, "repository.User.Delete")
//	}
//	log.Infof("Mongo Delete %v", res)
//	return nil
//}
//
//func (r *userDBRepo) Update(filter map[string]interface{}, fields map[string]interface{}) error {
//	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
//	defer cancel()
//	collection := r.client.Database(r.database).Collection("users")
//	res, err := collection.UpdateOne(
//		ctx,
//		filter,
//		fields,
//	)
//	if err != nil {
//		return errors.Wrap(err, "repository.User.Update")
//	}
//	log.Infof("Mongo Update %v", res)
//	return nil
//}
//
//func (r *userDBRepo) EncryptPassword(password string) (string, error) {
//	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	if err != nil {
//		return "", err
//	}
//	return string(hashedPass), nil
//}
//
//func (r *userDBRepo) DecryptPassword(hashedPass, password string) error {
//	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
//	if err != nil {
//		return err
//	}
//	return nil
//}
