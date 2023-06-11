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

type Operator struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Online    int    `json:"online"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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

type Dialog struct {
	ID          int         `json:"id"`
	State       string      `json:"state"`
	Begin       string      `json:"begin"`
	End         string      `json:"end"`
	OperatorID  int         `json:"operator_id"`
	LastMessage LastMessage `json:"last_message"`
}

type LastMessage struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Created  string `json:"created"`
	ClientID int    `json:"client_id"`
}

// GetAllDialogsOpenByOperatorID returns all dialogs open by operator id
func (s *Chat2DeskService) GetAllDialogsOpenByOperatorID(operatorID int) []Dialog {
	var offser int
	var index int

	var dialogs []Dialog
	for {
		index++
		var totalDialogsCounted int = 0
		queryString := fmt.Sprintf("state=open&limit=200&order=desc&offset=%d", offser)
		path := fmt.Sprintf("%s/dialogs?operator_id=%s&%s", url, strconv.Itoa(operatorID), queryString)
		resp, err := providers.MakeRquest(http.MethodGet, path, s.Token, nil)
		if err != nil {
			os.Exit(1)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var dialog ResponseDialogs
		json.Unmarshal(data, &dialog)

		totalDialogsCounted = dialog.Meta.Total
		dialogs = append(dialogs, dialog.Data...)

		if len(dialog.Data) >= totalDialogsCounted {
			break
		}
	}
	filteredDialogs := verifyDatetimeDialogIsToday(dialogs)

	return filteredDialogs
}

// GetDialogsByOperator returns all dialogs opened by operator id
func (s *Chat2DeskService) GetDialogsByOperator(operatorID int) []Dialog {
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

	var dialog ResponseDialogs
	json.Unmarshal(data, &dialog)
	return dialog.Data
}

type ResponseMessages struct {
	Data []Message `json:"data"`
}

type Message struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Created    string `json:"created"`
	DialogID   int    `json:"dialog_id"`
	OperatorID int    `json:"operator"`
	ClientID   int    `json:"client_id"`
}

// GetMessageByDialogID returns all messages from dialog id
func (s *Chat2DeskService) GetMessageByDialogID(dialogID int) []Message {
	path := fmt.Sprintf("%s/messages?dialog_id=%s&limit=200&order=desc", url, strconv.Itoa(dialogID))
	resp, err := providers.MakeRquest(http.MethodGet, path, s.Token, nil)
	if err != nil {
		os.Exit(1)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var messages ResponseMessages
	json.Unmarshal(data, &messages)
	filteredMessage := verifyDatetimeMessagesIsToday(messages.Data)
	return filteredMessage
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

// verifyDatetimeMessagesIsToday filters messages that are not from today
func verifyDatetimeMessagesIsToday(messages []Message) []Message {
	today := time.Now().Format("2006-01-02")

	filteredMessage := []Message{}
	for _, m := range messages {
		messageIsToday := m.Created[0:10]
		if messageIsToday == today {
			filteredMessage = append(filteredMessage, m)
		}
	}

	return filteredMessage
}
