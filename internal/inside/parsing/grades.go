package parsing

import (
	"encoding/xml"
	"regexp"

	"github.com/mamaart/go-learn/internal/inside/models"
)

type gradesTable struct {
	Rows []struct {
		Cols []struct {
			Cell string `xml:",innerxml"`
		} `xml:"td"`
	} `xml:"tr"`
}

func ParseGradesHtml(data []byte) (results []models.Grade) {
	pat := `<table[^>]*class="gradesList"[^>]*>[\s\S]*?<\/table>`
	var x gradesTable
	xml.Unmarshal(regexp.MustCompile(pat).Find(data), &x)
	for _, r := range x.Rows[1:] {
		results = append(results, models.Grade{
			URL:      extractURL(string(r.Cols[0].Cell)),
			Title:    string(r.Cols[1].Cell),
			Grade:    firstDigit(string(r.Cols[2].Cell)),
			Ects:     firstDigit(string(r.Cols[3].Cell)),
			Semester: string(r.Cols[4].Cell),
		})
	}
	return results
}

func extractURL(cell string) string {
	re := regexp.MustCompile(`<a href="([^"]+)"`)
	match := re.FindStringSubmatch(cell)
	if len(match) == 2 {
		return match[1]
	}
	return ""
}

func firstDigit(cell string) string {
	match := regexp.MustCompile(`(\d+)`).FindStringSubmatch(cell)
	if len(match) > 0 {
		return match[0]
	}
	return ""
}
