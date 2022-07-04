package qa

import (
	"fmt"
	"testing"
)

func Test_qa(t *testing.T) {
	qaDebug.Set("name", "123")
	fmt.Println(qaDebug.GetOptions())
}
