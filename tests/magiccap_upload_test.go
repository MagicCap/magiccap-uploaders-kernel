package tests

import (
	"encoding/json"
	"io/ioutil"
	"github.com/magiccap/magiccap-uploaders-kernel"
	"github.com/magiccap/magiccap-uploaders-kernel/standards"
	"os"
	"testing"
)

func TestMagicCapUpload(t *testing.T) {
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
	Uploader := k.Uploaders["magiccap"]
	if Uploader == nil {
		t.Error("i.magiccap is nil.")
	}
	if os.Getenv("MAGICCAP_INSTALL_ID") == "" {
		t.Skip("MAGICCAP_INSTALL_ID is not set. Skipping test!")
		return
	}
	map_ := map[string]interface{}{
		"install_id": os.Getenv("MAGICCAP_INSTALL_ID"),
	}
	url, err := k.Uploaders["magiccap"].Upload(map_, b, "magiccap.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
}
