package tests

import (
	"encoding/json"
	"io/ioutil"
	"magiccap-uploaders-kernel"
	"magiccap-uploaders-kernel/standards"
	"os"
	"testing"
)

func TestElixireUpload(t *testing.T) {
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
	b, err := ioutil.ReadFile("magiccap.png")
	if err != nil {
		t.Error(err)
	}
	Uploader := k.Uploaders["elixire"]
	if Uploader == nil {
		t.Error("Elixi.re is nil.")
	}
	if os.Getenv("ELIXIRE_TOKEN") == "" {
		t.Skip("ELIXIRE_TOKEN is not set. Skipping test!")
		return
	}
	map_ := map[string]interface{}{
		"elixire_token": os.Getenv("ELIXIRE_TOKEN"),
	}
	url, err := k.Uploaders["elixire"].Upload(map_, b, "magiccap.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
}
