package parsing

import (
	"encoding/xml"
	"errors"
	"io"
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

func getTable(data []byte) (table gradesTable, err error) {
	pat := `<table[^>]*class="gradesList"[^>]*>[\s\S]*?<\/table>`
	data = regexp.MustCompile(pat).Find(data)
	if len(data) == 0 {
		return table, errors.New("gradesList not found")
	}
	return table, xml.Unmarshal(data, &table)

}

func ParseGradesHtml(data []byte) (results []models.Grade, err error) {
	table, err := getTable(data)
	if err != nil {
		if err == io.EOF {
			return results, nil
		}
		return results, err
	}
	for _, row := range table.Rows[1:] {
		results = append(results, models.Grade{
			URL:      extractURL(string(row.Cols[0].Cell)),
			Title:    string(row.Cols[1].Cell),
			Grade:    firstDigit(string(row.Cols[2].Cell)),
			Ects:     firstDigit(string(row.Cols[3].Cell)),
			Semester: string(row.Cols[4].Cell),
		})
	}
	return results, nil
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
