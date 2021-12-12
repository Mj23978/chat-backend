package message

// type MessageService interface {
// 	FindById(id int) (*m.Message, error)
// 	Create(message *m.Message) error
// }

// type messageService struct {
// 	messageRepo MessageRepository
// }

// func NewMessageService(messageRepo MessageRepository) MessageService {
// 	return &messageService{
// 		messageRepo,
// 	}
// }

// func (svc *messageService) Create(message *m.Message) error {
// 	if e := validate.Validate(message); e != nil {
// 		return errs.Wrap(e, "service.Message.Store")
// 	}
// 	// message.Password = repo.EncryptPassword(data.Password)
// 	return svc.messageRepo.Create(message)
// }

// func (svc *messageService) FindById(id int) (*m.Message, error) {
// 	filter := map[string]interface{}{"id": id}
// 	res, e := svc.messageRepo.FindBy(filter)
// 	if e != nil {
// 		return res, errs.Wrap(e, "service.News.GetById")
// 	}
// 	return res, nil
// }
