package lib

import (
	"testing"
)

func Test_SnowGenerate(t *testing.T) {
	if Generate() <= 0 {
		t.Fatal("")
	}
}
