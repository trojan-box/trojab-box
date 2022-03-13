package util

import (
	"fmt"
	"testing"
)

func TestGetRandom(t *testing.T) {

	for i := 0; i < 100; i++ {
		result := GetRandom(0, 46)
		fmt.Printf("%d:%d\n", i, result)
	}
}
