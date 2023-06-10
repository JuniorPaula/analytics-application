package services

import (
	"c2d-reports/internal/providers"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var url = "https://api.chat24.io/v1"

type Chat2DeskService struct {
	Token string
}

type ResponseOperators struct {
	Data []Operator `json:"data"`
}

type Operator struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Online    int    `json:"online"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *Chat2DeskService) GetOperators() []Operator {
	path := fmt.Sprintf("%s/operators?limit=200", url)
	resp, err := providers.MakeRquest(http.MethodGet, path, s.Token, nil)

	if err != nil {
		os.Exit(1)
	}

	op, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var operators ResponseOperators
	json.Unmarshal(op, &operators)

	validOperators := formatOperators(operators)

	return validOperators
}

func formatOperators(op ResponseOperators) []Operator {
	var operators []Operator
	for i := 0; i < len(op.Data); i++ {
		if op.Data[i].Role != "deleted" && op.Data[i].Role != "disabled" {
			operators = append(operators, op.Data[i])
		}
	}

	return operators
}
