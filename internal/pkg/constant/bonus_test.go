package constant

import (
	"fmt"
	"github.com/thoas/go-funk"
	"testing"
)

func TestContainsBonus(t *testing.T) {
	bonusLevels := []BonusLevel{
		BonusLevel1,
		BonusLevel2,
	}
	fmt.Println(funk.Contains(bonusLevels, BonusLevel1))
	fmt.Println(funk.Contains(bonusLevels, BonusLevel3))
}
