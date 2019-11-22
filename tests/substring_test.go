package tests

import (
	"github.com/magiccap/magiccap-uploaders-kernel/utils"
	"testing"
)

func TestSubString(t *testing.T) {
	map_ := map[string]interface{}{
		"hello": "Hello",
		"world": "world",
	}
	subbed, err := utils.SubString("[SUB] {hello} {world}!", map_, "")
	if err != nil {
		t.Error(err.Error())
	} else {
		cmp := "[SUB] Hello world!"
		if subbed == cmp {
			t.Logf("Substitution test passed (\"%s\" == \"%s\")!", subbed, cmp)
		} else {
			t.Errorf("Substitution test failed (\"%s\" != \"%s\")!", subbed, cmp)
		}
	}
}
