package model

type BonusRecord struct {
	Model
	Address   string           `json:"address" gorm:"type:varchar(50)"`
	Bonus     int64            `json:"bonus"`
	Type      BonusRecordType  `json:"type"`
	Associate int64            `json:"associate"`
	State     BonusRecordState `json:"state"`
}

func (b *BonusRecord) TableName() string {
	return "t_bonus_record"
}

type BonusRecordType int

const (
	BonusRecordTypeWin      BonusRecordType = 1
	BonusRecordTypeWithdraw BonusRecordType = 2
	BonusRecordTypeShare    BonusRecordType = 3
)

type BonusRecordState int

const (
	BonusRecordStateSubmit     BonusRecordState = 0
	BonusRecordStateProcessing BonusRecordState = 1
	BonusRecordStateCompleted  BonusRecordState = 2
)
