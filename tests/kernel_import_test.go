package tests

import (
	"encoding/json"
	"io/ioutil"
	"magiccap-uploaders-kernel"
	"magiccap-uploaders-kernel/standards"
	"testing"
)

func TestV1Import(t *testing.T) {
	k := MagicCapKernel.Kernel{
		Uploaders: map[string]*MagicCapKernelStandards.Uploader{},
	}
	f, err := ioutil.ReadFile("../uploaders/v1.json")
	if err != nil {
		t.Error("Failed to open V1 uploaders. Does it exist?")
	}
	var j map[string]interface{}
	err = json.Unmarshal(f, &j)
	if err != nil {
		t.Error(err)
	}
	err = k.Load(j)
	if err != nil {
		t.Error(err)
	}
	t.Log("Everything imported fine!")
}
