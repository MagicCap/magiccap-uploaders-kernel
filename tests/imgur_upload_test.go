package tests

import (
	"encoding/json"
	"io/ioutil"
	"magiccap-uploaders-kernel"
	"magiccap-uploaders-kernel/standards"
	"testing"
)

func TestImgurUpload(t *testing.T) {
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
	Uploader := k.Uploaders["imgur"]
	if Uploader == nil {
		t.Error("Imgur is nil.")
	}
	url, err := k.Uploaders["imgur"].Upload(make(map[string]interface{}, 0), b, "magiccap.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
}
