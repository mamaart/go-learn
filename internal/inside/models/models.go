package models

import "encoding/json"

type Grade struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Grade    string `json:"grade"`
	Ects     string `json:"ects"`
	Semester string `json:"semester"`
}

func (m Grade) String() string {
	jsond, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "fail"
	}
	return string(jsond)
}
