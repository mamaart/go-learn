package models

import "encoding/json"

type Grade struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Grade    int    `json:"grade"`
	Ects     int    `json:"ects"`
	Semester string `json:"semester"`
}

func (m Grade) String() string {
	jsond, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "fail"
	}
	return string(jsond)
}
