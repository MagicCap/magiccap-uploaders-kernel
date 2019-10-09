package MagicCapKernelStandards

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/ajg/form.v1"
	"io/ioutil"
	"magiccap-uploaders-kernel/utils"
	"math"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

// POSTAs defines what the type of the POST request is.
type POSTAs struct {
	Type string `json:"type"`
	Key string `json:"key"`
}

// HTTPSpec defines the HTTP spec for this uploader.
type HTTPSpec struct {
	Method string `json:"method"`
	URL string `json:"url"`
	POSTAs *POSTAs `json:"post_as"`
	Headers *map[string]string `json:"headers"`
	Response *string `json:"response"`
}

// HTTPInit defines the HTTP standard.
func HTTPInit(Structure UploaderStructure) (*Uploader, error) {
	b, err := json.Marshal(Structure.Spec)
	if err != nil {
		return nil, err
	}
	var spec HTTPSpec
	err = json.Unmarshal(b, &spec)
	if err != nil {
		return nil, err
	}
	e := base64.Encoding{}
	return &Uploader{
		Description:   Structure.Description,
		Name:          Structure.Name,
		ConfigOptions: Structure.Config,
		Icon:          Structure.Icon,
		Upload: func(Config map[string]interface{}, Data []byte, Filename string) (string, error) {
			var URL string
			var POSTData *bytes.Buffer
			var ContentType *string

			if spec.POSTAs.Type == "b64" {
				URL = spec.URL + "?" + spec.POSTAs.Key + "=" + e.EncodeToString(Data)
			} else if spec.POSTAs.Type == "raw" {
				POSTData = bytes.NewBuffer(Data)
				c := http.DetectContentType(Data)
				ContentType = &c
			} else if spec.POSTAs.Type == "multipart" {
				buffer := new(bytes.Buffer)
				writer := multipart.NewWriter(buffer)
				if err != nil {
					return "", err
				}
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, spec.POSTAs.Key, Filename))
				h.Set("Content-Type", http.DetectContentType(Data))
				part, err := writer.CreatePart(h)
				if err != nil {
					return "", err
				}
				_, err = part.Write(Data)
				if err != nil {
					return "", err
				}
				err = writer.Close()
				if err != nil {
					return "", err
				}
				POSTData = buffer
				URL = spec.URL
				c := writer.FormDataContentType()
				ContentType = &c
			} else if spec.POSTAs.Type == "urlencoded" {
				u, err := form.EncodeToString(map[string]interface{}{
					spec.POSTAs.Key: Data,
				})
				if err != nil {
					return "", err
				}
				POSTData = bytes.NewBufferString(u)
				URL = spec.URL
				c := "application/x-www-form-urlencoded"
				ContentType = &c
			} else {
				return "", errors.New("POST type not defined.")
			}
			URL, err = utils.SubString(URL, Config, Filename)
			if err != nil {
				return "", err
			}

			r, err := http.NewRequest(spec.Method, URL, POSTData)
			if err != nil {
				return "", err
			}
			if spec.Headers != nil {
				for k, v := range *spec.Headers {
					v, err = utils.SubString(v, Config, Filename)
					if err != nil {
						return "", err
					}
					r.Header.Set(k, v)
				}
			}
			if ContentType != nil {
				r.Header.Set("Content-Type", *ContentType)
			}
			client := http.Client{
				Timeout: 30 * time.Second,
			}
			resp, err := client.Do(r)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", nil
			}
			ResponseType := math.Floor(float64(resp.StatusCode) / 100)
			if ResponseType == 4 || ResponseType == 5 {
				return "", errors.New("Uploader returned the status " + strconv.Itoa(resp.StatusCode) + ".")
			}
			if spec.Response == nil {
				return string(b), nil
			} else {
				var JSONMap map[string]interface{}
				err := json.Unmarshal(b, &JSONMap)
				if err != nil {
					return "", err
				}
				FinalURL := *spec.Response
				for true {
					full := ""
					sub := ""
					for _, char := range FinalURL {
						if full != "" {
							if char == '%' {
								full += "%"
								break
							} else {
								full += string(char)
								sub += string(char)
							}
						} else if char == '%' {
							full = "%"
						}
					}
					if full == "" {
						break
					}
					Key := strings.Split(sub, ".")
					MapContext := JSONMap
					ArrContext := make([]interface{}, 0)
					IsMap := true
					Last, Key := Key[len(Key)-1], Key[:len(Key)-1]
					var ok bool
					for _, v := range Key {
						if IsMap {
							NewItem := MapContext[v]
							MapContext, ok = NewItem.(map[string]interface{})
							if !ok {
								ArrContext = NewItem.([]interface{})
								IsMap = false
							}
						} else {
							// This is an array.
							i, err := strconv.Atoi(v)
							if err != nil {
								return "", err
							}
							Item := ArrContext[i]
							MapContext, ok = Item.(map[string]interface{})
							if !ok {
								// Is this an array?
								ArrContext, ok = Item.([]interface{})
								if !ok {
									return "", errors.New("Failed to parse as a map or an array.")
								}
								IsMap = false
							} else {
								IsMap = true
							}
						}
					}
					var s string
					if IsMap {
						s, ok = MapContext[Last].(string)
						if !ok {
							return "", errors.New("The final value in the uploader is not a string.")
						}
					} else {
						i, err := strconv.Atoi(Last)
						if err != nil {
							return "", err
						}
						s, ok = ArrContext[i].(string)
						if !ok {
							return "", errors.New("The final value in the uploader is not a string.")
						}
					}
					FinalURL = strings.Replace(FinalURL, full, s, 1)
					full = ""
					sub = ""
				}
				FinalURL, err = utils.SubString(FinalURL, Config, Filename)
				if err != nil {
					return "", err
				}
				return FinalURL, nil
			}
		},
	}, nil
}
