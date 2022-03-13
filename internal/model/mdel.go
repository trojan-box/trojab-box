package model

type Model struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt int64
	UpdatedAt int64
}
