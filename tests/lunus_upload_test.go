package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"github.com/magiccap/magiccap-uploaders-kernel"
	"github.com/magiccap/magiccap-uploaders-kernel/standards"
	"math"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestLunusUpload(t *testing.T) {
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
	Uploader := k.Uploaders["novus"]
	if Uploader == nil {
		t.Error("Lunus is nil.")
	}
	Token := os.Getenv("NOVUS_TOKEN")
	if Token == "" {
		t.Skip("NOVUS_TOKEN is not set. Skipping test!")
		return
	}
	map_ := map[string]interface{}{
		"novus_token": Token,
	}
	url, err := k.Uploaders["novus"].Upload(map_, b, "magiccap.png")
	if err != nil {
		t.Error(err)
		return
	}
	res, err := http.NewRequest("DELETE", url, bytes.NewBuffer([]byte{}))
	res.Header.Set("Authorization", "Bearer " + Token)
	if err != nil {
		t.Error(err)
		return
	}
	client := http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		t.Error(err)
		return
	}
	status := math.Floor(float64(resp.StatusCode) / 100)
	if status == 4 || status == 5 {
		t.Error("Deletion returned the status " + strconv.Itoa(resp.StatusCode) + ".")
		return
	}
	t.Log("Successfully uploaded and deleted the image (per owners request).")
}
