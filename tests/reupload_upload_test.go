package tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	MagicCapKernel "github.com/magiccap/magiccap-uploaders-kernel"
	MagicCapKernelStandards "github.com/magiccap/magiccap-uploaders-kernel/standards"
)

func TestReuploadUpload(t *testing.T) {
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
	Uploader := k.Uploaders["reupload"]
	if Uploader == nil {
		t.Error("reupload is nil.")
	}
	if os.Getenv("REUPLOAD_TOKEN") == "" {
		t.Skip("REUPLOAD_TOKEN is not set. Skipping test!")
		return
	}
	map_ := map[string]interface{}{
		"reupload_token": os.Getenv("REUPLOAD_TOKEN"),
	}
	url, err := k.Uploaders["reupload"].Upload(map_, b, "magiccap.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
}
