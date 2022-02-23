package check_tool

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type ValidFunc func(atom Atom, args string) error
type MapFunc map[string]ValidFunc

func numberFunc(atom Atom, _ string) error {
	valid, err := regexp.MatchString("^[0-9]+$", atom.Value)
	if err != nil {
		log.Printf("ERROR: check_tool.Number: %v\n", err)
		return ErrorKCHECK
	}

	if valid {
		return nil
	}

	message := "all characters must be numeric, invalid value `%s` in field `%s`"
	return fmt.Errorf(message, atom.Value, atom.Name)
}

func decimalFunc(atom Atom, _ string) error {
	valid, err := regexp.MatchString("^[0-9]+.[0-9]+$", atom.Value)
	if err != nil {
		log.Printf("ERROR: check_tool.Number: %v\n", err)
		return ErrorKCHECK
	}

	if valid {
		return nil
	}

	message := "decimal expected, invalid value `%s` in field `%s`"
	return fmt.Errorf(message, atom.Value, atom.Name)
}

func sword(atom Atom, _ string) error {
	valid, err := regexp.MatchString("^[0-9a-zA-Z_ñ]*$", atom.Value)
	if err != nil {
		log.Printf("ERROR: check_tool.sword: %v\n", err)
		return ErrorKCHECK
	}

	if valid {
		return nil
	}

	message := "in field `%s` only numeric and alphabetic characters are allowed"
	return fmt.Errorf(message, atom.Name)
}

// calLens returns the stringLen value converted to int, the number of characters in the value, and error if it exists.
// Used by Length, MaxLength, MinLength
func calLens(value string, stringLen string) (int, int, error) {
	length, err := strconv.Atoi(stringLen)
	if err != nil {
		log.Printf("ERROR: check_tool.calLens: %v\n", err)
		return 0, 0, ErrorKCHECK
	}

	valueLength := len(value)
	return length, valueLength, nil
}

func noNilFunc(atom Atom, _ string) error {
	length := len(atom.Value)
	if strings.TrimSpace(atom.Value) == "" {
		message := "field `%s` cannot be empty"
		if length != 0 {
			message = "field `%s` cannot contain whitespace only"
		}

		return fmt.Errorf(message, atom.Name)
	}

	return nil
}

func noSpacesStartAndEnd(atom Atom, _ string) error {
	matchStartSpace, _ := regexp.MatchString("^( .)", atom.Value)
	if matchStartSpace {
		message := "field `%s` cannot contain leading spaces"
		return fmt.Errorf(message, atom.Name)
	}

	matchEndSpace, _ := regexp.MatchString("(. )$", atom.Value)
	if matchEndSpace {
		message := "field `%s` cannot contain trailing spaces"
		return fmt.Errorf(message, atom.Name)
	}

	return nil
}

func textFunc(atom Atom, args string) error {
	denied := "!\"#$%&'()*+,./:;<=>?@[\\]^_}{~|"
	if err := noSpacesStartAndEnd(atom, args); err != nil {
		return err
	}

	match, _ := regexp.MatchString("( ){3}", atom.Value)
	if match {
		const message = "field `%s` cannot have words separated by more than 2 spaces"
		return fmt.Errorf(message, atom.Name)
	}

	for _, c := range atom.Value {
		if strings.ContainsRune(denied, c) {
			const message = "field `%s` cannot contain any of these characters %s"
			return fmt.Errorf(message, atom.Name, denied)

		}
	}

	return nil
}

func emailFunc(atom Atom, _ string) error {
	match, err := regexp.MatchString(`^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$`, atom.Value)
	if err != nil {
		log.Printf("ERROR: check_tool.emailFunc: %v\n", err)
		return ErrorKCHECK
	}

	if !match {
		message := "el campo `%s` es del tipo correo, `%s` no es un correo válido"
		return fmt.Errorf(message, atom.Name, atom.Value)
	}

	return nil
}

func lengthFunc(atom Atom, args string) error {
	vLen, valueLength, err := calLens(atom.Value, args)
	if err != nil {
		log.Printf("ERROR: check_tool.lenghtFunc field:`%s` args:`%s`: %v\n", atom.Name, args, err)
		return err
	}

	if valueLength == vLen {
		return nil
	}

	message := "number of characters in field `%s` must be `%d`, `%s` has `%d` characters"
	return fmt.Errorf(message, atom.Name, vLen, atom.Value, valueLength)
}

func maxLengthFunc(atom Atom, args string) error {
	maxLen, valueLength, err := calLens(atom.Value, args)
	if err != nil {
		log.Printf("ERROR: check_tool.maxLenghtFunc field:`%s` args:`%s`: %v\n", atom.Name, args, err)
		return err
	}

	if valueLength <= maxLen {
		return nil
	}

	message := "field `%s` maximum number of characters must be `%d`, `%s` has `%d` characters"
	return fmt.Errorf(message, atom.Name, maxLen, atom.Value, valueLength)
}

func minLengthFunc(atom Atom, args string) error {
	minLen, valueLength, err := calLens(atom.Value, args)
	if err != nil {
		log.Printf("ERROR: check_tool.minLenghtFunc field:`%s` args:`%s`: %v\n", atom.Name, args, err)
		return err
	}

	if valueLength >= minLen {
		return nil
	}

	message := "the minimum number of characters in the `%s` field must be `%d`, `%s` has `%d` characters"
	return fmt.Errorf(message, atom.Name, minLen, atom.Value, valueLength)
}

func regularExpression(atom Atom, args string) error {
	valid, err := regexp.MatchString(args, atom.Value)
	if err != nil {
		log.Printf("ERROR: check_tool.regularExpression: %v, in field `%s` with expression `%s`\n", err, atom.Name, args)
		return ErrorKCHECK
	}

	if valid {
		return nil
	}

	message := "the value `%s` in the field `%s` is invalid, check with the administrator for more information"
	return fmt.Errorf(message, atom.Value, atom.Name)
}
