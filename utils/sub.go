package utils

import (
	"encoding/json"
	"strings"
)

func SubString(v string, Config map[string]interface{}, Filename string) (string, error) {
	for true {
		full := ""
		sub := ""
		for _, char := range v {
			if full != "" {
				if char == '}' {
					full += "}"
					break
				} else {
					full += string(char)
					sub += string(char)
				}
			} else if char == '{' {
				full = "{"
			}
		}
		if full == "" {
			break
		}
		c := Config[sub]
		Result, ok := c.(string)
		if !ok {
			if sub == "ext" {
				s := strings.Split(Filename, ".")
				Result = s[len(s)-1]
			} else {
				b, err := json.Marshal(&c)
				if err != nil {
					return "", err
				}
				Result = string(b)
			}
		}
		v = strings.Replace(v, full, Result, 1)
		full = ""
		sub = ""
	}
	return v, nil
}
