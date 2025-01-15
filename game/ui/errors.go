package ui

import (
	"strings"
)

// ErrorList представляет собой список ошибок, которые можно собирать и обрабатывать как одну ошибку.
type ErrorList struct {
	errors []error
}

// Add добавляет ошибку в список, если она не равна nil.
func (e *ErrorList) Add(err error) {
	if err != nil {
		e.errors = append(e.errors, err)
	}
}

// Error возвращает строковое представление всех ошибок в списке, объединенных через точку с запятой.
func (e *ErrorList) Error() string {
	if len(e.errors) == 0 {
		return ""
	}

	errStrings := make([]string, len(e.errors))
	for i, err := range e.errors {
		errStrings[i] = err.Error()
	}
	return strings.Join(errStrings, "; ")
}

// HasErrors возвращает true, если в списке есть хотя бы одна ошибка.
func (e *ErrorList) HasErrors() bool {
	return len(e.errors) > 0
}

// Errors возвращает слайс всех ошибок в списке.
func (e *ErrorList) Errors() []error {
	return e.errors
}

// Clear очищает список ошибок.
func (e *ErrorList) Clear() {
	e.errors = nil
}
