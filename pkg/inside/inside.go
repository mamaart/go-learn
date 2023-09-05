package inside

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mamaart/go-learn/internal/auth"
	"github.com/mamaart/go-learn/internal/inside/models"
	"github.com/mamaart/go-learn/internal/inside/parsing"
)

type Inside struct {
	cli *http.Client
}

func New(username, password string) (*Inside, error) {
	cli, err := auth.LoginToInside(username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to login %s", err)
	}

	i := Inside{cli}

	return &i, nil

}

func (i *Inside) GetGrades() ([]models.Grade, error) {
	url := "https://cn.inside.dtu.dk/cnnet/Grades/Grades.aspx"
	resp, err := i.cli.Get(url)
	if err != nil {
		return nil, fmt.Errorf("data parse failed: %s", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parsing.ParseGradesHtml(data), nil
}
