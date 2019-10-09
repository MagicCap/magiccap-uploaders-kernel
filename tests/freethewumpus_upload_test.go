package tests

import (
	"encoding/json"
	"io/ioutil"
	"magiccap-uploaders-kernel"
	"magiccap-uploaders-kernel/standards"
	"os"
	"testing"
)

func TestFreeTheWumpusUpload(t *testing.T) {
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
	Uploader := k.Uploaders["freethewumpus"]
	if Uploader == nil {
		t.Error("freethewump.us is nil.")
	}
	if os.Getenv("FTW_TOKEN") == "" {
		t.Skip("FTW_TOKEN is not set. Skipping test!")
		return
	}
	map_ := map[string]interface{}{
		"ftw_token": os.Getenv("FTW_TOKEN"),
	}
	url, err := k.Uploaders["freethewumpus"].Upload(map_, b, "magiccap.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
}
