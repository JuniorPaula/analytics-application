package services

import (
	"c2d-reports/internal/providers"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var url = "https://api.chat24.io/v1"

type Chat2DeskService struct {
	Token string
}

type ResponseOperators struct {
	Data []Operator `json:"data"`
}

// Operator is the model for the operators
type Operator struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	Online        int    `json:"online"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	OpenedDialogs int    `json:"opened_dialogs"`
}

// GetOperators returns all operators from company
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

// formatOperators filters operators that are deleted or disabled
func formatOperators(op ResponseOperators) []Operator {
	var operators []Operator
	for i := 0; i < len(op.Data); i++ {
		if op.Data[i].Role != "deleted" && op.Data[i].Role != "disabled" {
			operators = append(operators, op.Data[i])
		}
	}

	return operators
}

type ResponseDialogs struct {
	Data []Dialog `json:"data"`
	Meta Meta     `json:"meta"`
}

type Meta struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
}

type ResponseDialog struct {
	Data Dialog `json:"data"`
}

// Dialog is a struct that represents a dialog
type Dialog struct {
	ID          int         `json:"id"`
	State       string      `json:"state"`
	Begin       string      `json:"begin"`
	End         string      `json:"end"`
	OperatorID  int         `json:"operator_id"`
	LastMessage LastMessage `json:"last_message"`
}

type LastMessage struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Created   string `json:"created"`
	ClientID  int    `json:"client_id"`
	RequestID int    `json:"request_id"`
}

// GetAllDialogsOpenByOperatorID returns all dialogs open by operator id
func (s *Chat2DeskService) GetAllDialogsOpenByOperatorID(operatorID int) []Dialog {
	queryString := "state=open&limit=200&order=desc"
	path := fmt.Sprintf("%s/dialogs?operator_id=%s&%s", url, strconv.Itoa(operatorID), queryString)
	resp, err := providers.MakeRquest(http.MethodGet, path, s.Token, nil)
	if err != nil {
		os.Exit(1)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var dialogs ResponseDialogs
	json.Unmarshal(data, &dialogs)

	filteredDialogs := verifyDatetimeDialogIsToday(dialogs.Data)

	return filteredDialogs
}

// verifyDatetimeDialogIsToday filters dialogs that are not from today
// and returns only dialogs from today
func (s *Chat2DeskService) GetDialogByID(dialogID int) Dialog {
	path := fmt.Sprintf("%s/dialogs/%s", url, strconv.Itoa(dialogID))
	resp, err := providers.MakeRquest(http.MethodGet, path, s.Token, nil)
	if err != nil {
		os.Exit(1)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var dialog ResponseDialog
	json.Unmarshal(data, &dialog)

	return dialog.Data
}

type ResposeClients struct {
	Data Client `json:"data"`
}

// Client represents a client from chat2desk
type Client struct {
	ID           int                      `json:"id"`
	Name         string                   `json:"name"`
	AssignedName string                   `json:"assigned_name"`
	Phone        string                   `json:"phone"`
	Tags         []struct{ Label string } `json:"tags"`
}

// GetClientByID returns client by id
func (s *Chat2DeskService) GetClientByID(clientID int) Client {
	path := fmt.Sprintf("%s/clients/%s", url, strconv.Itoa(clientID))
	resp, err := providers.MakeRquest(http.MethodGet, path, s.Token, nil)
	if err != nil {
		os.Exit(1)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var client ResposeClients
	json.Unmarshal(data, &client)

	return client.Data
}

// verifyDatetimeDialogIsToday filters dialogs that are not from today
func verifyDatetimeDialogIsToday(dialogs []Dialog) []Dialog {
	today := time.Now().Format("2006-01-02")

	filteredDialogs := []Dialog{}
	for _, dialog := range dialogs {
		dialogIsToday := dialog.LastMessage.Created[0:10]
		if dialogIsToday == today {
			filteredDialogs = append(filteredDialogs, dialog)
		}
	}

	return filteredDialogs
}
