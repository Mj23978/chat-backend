package chat

type CallLogService interface {
	Get(id string) error
	Delete(id string) error
	Update(id string, fields map[string]interface{}) error
	Create(id string) error
	// Apply Filter Parameter for Costume Gets to All GetAll Functions
	GetAll() error
}

type MessageService interface {
	Get(id string) error
	Delete(id string) error
	Update(id string, fields map[string]interface{}) error
	Create(id string) error
	Send() error
	UploadFile()
}

type GroupService interface {
	Get(id string) error
	Delete(id string) error
	Update(id string, fields map[string]interface{}) error
	Create(id string) error
	GetAll() error
}

type StoryService interface {
	Get(id string) error
	Delete(id string) error
	Update(id string, fields map[string]interface{}) error
	Create(id string) error
	GetAll() error
	GetSubscribed()
}

type CommentService interface {
	Get(id string) error
	Delete(id string) error
	Update(id string, fields map[string]interface{}) error
	Create(id string) error
	GetAll() error
	GetMessage() error
}

type TagService interface {
	Get(id string) error
	Delete(id string) error
	Update(id string, fields map[string]interface{}) error
	Create(id string) error
	GetAll() error
	GetRelatedMessages()
}
