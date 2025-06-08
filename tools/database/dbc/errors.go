package dbc

import "fmt"

type ParseError struct {
	Source string
	Field  string
	Reason string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parse error in %s (field: %s): %s", e.Source, e.Field, e.Reason)
}

type ValidationError struct {
	Type  string
	ID    int
	Field string
	Msg   string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s ID %d (field: %s): %s", e.Type, e.ID, e.Field, e.Msg)
}

type DataLoadError struct {
	Source   string
	DataType string
	Reason   string
}

func (e DataLoadError) Error() string {
	return fmt.Sprintf("failed to load %s from %s: %s", e.DataType, e.Source, e.Reason)
}
