package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Operation struct {
	// Название выполненой операции. Как правильно
	// это текст соответсвующего http-кода ответа.
	Title string `json:"title"`
	// http-код операции.
	Status int `json:"status"`
	// Человекопонятное текстовое описание
	// выполненной операции
	Detail string `json:"detail"`
}

// Описание успешно выполненной операции
type SuccessOperation struct {
	Operation
	// Результат выполнения операции
	Result interface{} `json:"result"`
}

// Описание проблемы возникшей в работе API
// За основу взят стандарт https://tools.ietf.org/html/rfc7807
type ProblemOperation struct {
	Operation
	// Содержит URL с описанием проблем данного типа.
	Type string `json:"type"`
	// Содержит URL с подробным описанием конкретной проблемы.
	Instance string `json:"instance"`
}

func sendSuccess(status int, detail string, args ...interface{}) interface{} {
	return SuccessOperation{
		Operation: Operation{
			Title:  http.StatusText(status),
			Status: status,
			Detail: fmt.Sprintf(detail, args...),
		},
	}
}

func sendSuccessOK(detail string, args ...interface{}) []byte {
	response := sendSuccess(http.StatusOK, detail, args...)
	body, _ := json.Marshal(response)
	return body
}

func sendSuccessOKWithResult(detail string, result interface{}) []byte {
	response := SuccessOperation{
		Operation: Operation{
			Title:  http.StatusText(http.StatusOK),
			Status: http.StatusOK,
			Detail: detail,
		},
		Result: result,
	}
	body, _ := json.Marshal(response)
	return body
}

func sendProblem(status int, detail string, args ...interface{}) interface{} {
	return ProblemOperation{
		Operation: Operation{
			Title:  http.StatusText(status),
			Status: status,
			Detail: fmt.Sprintf(detail, args...),
		},
	}
}

func sendProblemNotFound(detail string, args ...interface{}) (int, interface{}) {
	return 0, sendProblem(http.StatusNotFound, detail, args...)
}

func sendProblemBadRequest(detail string, args ...interface{}) []byte {
	response := sendProblem(http.StatusBadRequest, detail, args...)
	body, _ := json.Marshal(response)
	return body
}

func sendProblemInternalServerError(detail string, args ...interface{}) (int, interface{}) {
	return 0, sendProblem(http.StatusInternalServerError, detail, args...)
}

func sendProblemPageNotFound(id int) (int, interface{}) {
	return sendProblemNotFound("Page \"%d\" not found.", id)
}

func sendProblemSpecNotFound(id int) (int, interface{}) {
	return sendProblemNotFound("Spec \"%d\" not found.", id)
}

func sendProblemParameterMustBeInteger(parameterName string, err error) (int, interface{}) {
	return 0, sendProblemBadRequest(
		errors.Wrapf(err, "Can't parse parameter \"%s\". "+
			"Parameter must be integer", parameterName).Error(),
	)
}
