package model

type UploadState int32

const (
	Preparing UploadState = iota
	Inprogress
	Success
	Failed
)

type Attachment struct {
	ID         string `json:"id"  msgpack:"id"  validate:"empty=false"`
	Type       string `json:"type"  msgpack:"type"`
	Title      string `json:"title" msgpack:"title"`
	TitleLink  string `json:"title_link"  msgpack:"title_link"`
	ThumbUrl   string `json:"thumb_url" msgpack:"thumb_url"`
	ImageUrl   string `json:"image_url" msgpack:"image_url"`
	Text       string `json:"text" msgpack:"text"`
	PreText    string `json:"pre_text" msgpack:"pre_text"`
	Footer     string `json:"footer" msgpack:"footer"`
	AuthorName string `json:"author_name" msgpack:"author_name"`
	AuthorLink string `json:"author_link" msgpack:"author_link"`
	File       string `json:"file" msgpack:"file"`
}

type AttachmentFile struct {
	Size  string  `json:"size"  msgpack:"size"  validate:"empty=false"`
	Path  string  `json:"path"  msgpack:"path"`
	Name  string  `json:"name"  msgpack:"name"`
	Bytes []uint8 `json:"bytes"  msgpack:"bytes"`
}

type Action struct {
	Name  string `json:"name"  msgpack:"name" validate:"empty=false"`
	Style string `json:"style"  msgpack:"style"`
	Text  string `json:"text"  msgpack:"text"`
	Type  string `json:"type"  msgpack:"type"`
	Value string `json:"value"  msgpack:"value"`
}
