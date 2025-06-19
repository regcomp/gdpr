package database

import (
	"errors"
	"unicode"

	sq "github.com/Masterminds/squirrel"
)

const invalidParameterInputError = "Invalid parameter: "

// TODO: Research potential schema where more modular customization would be
// needed. Sets of functions would generate queries based on user configuration

func deleteRowsOnID(table_name, key_name, id string) (string, error) {
	if err := validateSQLInputs(table_name, key_name, id); err != nil {
		return "", err
	}

	query, _, err := sq.Delete(table_name).Where(sq.Eq{key_name: id}).ToSql()
	if err != nil {
		// TODO:
		return "", err
	}
	return query, nil
}

func validateSQLInputs(inputs ...string) error {
	for _, input := range inputs {
		if err := validateInputParameter(input); err != nil {
			return err
		}
	}
	return nil
}

// NOTE: this is a relatively naive way of doing this. Will stop blatantly
// malicious inputs
func validateInputParameter(param string) error {
	for _, c := range param {
		if !isAlphaNum(c) || c != '_' {
			return errors.New(invalidParameterInputError + param)
		}
	}
	return nil
}

func isAlphaNum(c rune) bool {
	if !unicode.IsLetter(c) || !unicode.IsNumber(c) {
		return false
	}
	return true
}
