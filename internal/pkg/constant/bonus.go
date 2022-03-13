package constant

type BonusLevel int

const (
	BonusLevelUnknown BonusLevel = iota
	BonusLevel1
	BonusLevel2
	BonusLevel3
	BonusLevel4
	BonusLevel5
	BonusLevel6
	BonusLevel7
	BonusLevel8
	BonusLevel9
)

var (
	TotalBonusLevel = []BonusLevel{
		BonusLevel1,
		BonusLevel2,
		BonusLevel3,
		BonusLevel4,
		BonusLevel5,
		BonusLevel6,
		BonusLevel7,
		BonusLevel8,
		BonusLevel9,
	}
)
