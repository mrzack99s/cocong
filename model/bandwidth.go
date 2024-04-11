package model

type Bandwidth struct {
	BaseModel

	Name          string
	DownloadLimit int64
	UploadLimit   int64
}
