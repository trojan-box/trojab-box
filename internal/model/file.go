package model

type UploadFile struct {
	Model
	Address  string `json:"address" gorm:"type:varchar(50)"`
	IpfsHash string `json:"ipfs_hash"`
	Link     string `json:"link"`
}

func (l *UploadFile) TableName() string {
	return "t_file"
}
