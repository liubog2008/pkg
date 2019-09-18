package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Error defines errors which can be used by http frontend
type Error struct {
	// Code is defined as http status code
	Code int `json:"-"`
	// Reason defines short reason of error
	Reason string `json:"reason"`
	// Message defines long description of this error
	Message string `json:"message"`
	// Data defines some necessary data which is needed by frontend
	Data map[string]interface{} `json:"data"`
}

// Error implements error interface
func (e *Error) Error() string {
	return e.Message
}

// Factory defines error factory, it produce a set of errors with same type
type Factory interface {
	// New returns an error with data
	New(args ...interface{}) error
}

type factory struct {
	code     int
	reason   string
	template *Template
}

// MustNewFactory returns a new factory of error
// If format cannot be parsed, program will be panic
func MustNewFactory(code int, reason string, format string) Factory {
	f, err := NewFactory(code, reason, format)
	if err != nil {
		panic(err)
	}
	return f
}

// NewFactory returns a new facotry
// If format cannot be parsed, an error will be returned
func NewFactory(code int, reason string, format string) (Factory, error) {
	f := &factory{
		code:     code,
		reason:   reason,
		template: &Template{},
	}
	// factory should be new
	if err := f.template.fromRaw(format); err != nil {
		return nil, err
	}
	return f, nil
}

// New returns a formatted error
func (f *factory) New(args ...interface{}) error {
	data := convert(f.template.varNames, args)
	return &Error{
		Code:    f.code,
		Reason:  f.reason,
		Message: f.template.render(args),
		Data:    data,
	}
}

// Template defines template of errors
// e.g.
//   This is a %{hello}
type Template struct {
	format   string
	varNames []string
}

func (t *Template) fromRaw(raw string) error {
	sb := strings.Builder{}
	// isSign means previous is % and previous two are not double %
	// isOpen means '%{' have not meet '}' yet
	isSign, isOpen := false, false
	varNameBuilder := strings.Builder{}

	for _, r := range raw {
		switch r {
		case '%':
			if isOpen {
				return errors.New("% in {} is not allowed")
			}
			sb.WriteByte('%')
			isSign = !isSign

		case '{':
			if isSign {
				isSign, isOpen = false, true
			} else {
				sb.WriteByte('{')
			}
		case '}':
			if isOpen {
				varName := strings.TrimSpace(varNameBuilder.String())
				if varName == "" {
					return errors.New("param name should not be empty")
				}
				isOpen = false
				t.varNames = append(t.varNames, varName)
				varNameBuilder.Reset()
				// use %v as default output format
				sb.WriteByte('v')
			} else {
				sb.WriteByte('}')
			}
		default:
			if isOpen {
				varNameBuilder.WriteRune(r)
			} else {
				sb.WriteRune(r)
			}
		}
	}
	t.format = sb.String()
	return nil
}

func (t *Template) render(args []interface{}) string {
	return fmt.Sprintf(t.format, args...)
}

func convert(varNames []string, args []interface{}) map[string]interface{} {
	data := map[string]interface{}{}

	i := 0
	for i < len(varNames) && i < len(args) {
		data[varNames[i]] = args[i]
		i++
	}
	for k := i; k < len(args); k++ {
		n := fmt.Sprintf("(EXTRA)%v", args[k])
		data[n] = args[k]
	}
	return data
}
