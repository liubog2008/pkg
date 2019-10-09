package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFactory(t *testing.T) {
	cases := []struct {
		description string
		code        int
		reason      string
		format      string
		template    *Template
	}{
		{
			description: "new factory successfully",
			code:        http.StatusBadRequest,
			reason:      "FormatError",
			format:      "Expected format is '[a-z]*', actual %{format}",
			template: &Template{
				format: "Expected format is '[a-z]*', actual %v",
				varNames: []string{
					"format",
				},
			},
		},
		{
			description: "two variable",
			code:        http.StatusBadRequest,
			reason:      "FormatError",
			format:      "%{var} and %{another}",
			template: &Template{
				format: "%v and %v",
				varNames: []string{
					"var",
					"another",
				},
			},
		},
		{
			description: "variable with spaces",
			code:        http.StatusBadRequest,
			reason:      "FormatError",
			format:      "%{ var } and %{another}",
			template: &Template{
				format: "%v and %v",
				varNames: []string{
					"var",
					"another",
				},
			},
		},
		{
			description: "double %%",
			code:        http.StatusBadRequest,
			reason:      "FormatError",
			format:      "%%{ var } and %{another}",
			template: &Template{
				format: "%%{ var } and %v",
				varNames: []string{
					"another",
				},
			},
		},
	}
	for _, c := range cases {
		f, err := NewFactory(c.code, c.reason, c.format)
		assert.Nil(t, err, c.description)

		nf, ok := f.(*factory)
		assert.True(t, ok, c.description)
		assert.Equal(t, c.code, nf.code, c.description, ": code should be equal")
		assert.Equal(t, c.reason, nf.reason, c.description, ": reason should be equal")
		assert.Equal(t, c.template.varNames, nf.template.varNames, c.description, ": varnames should be equal")
		assert.Equal(t, c.template.format, nf.template.format, c.description, ": format should be equal")
	}
}

func TestNew(t *testing.T) {
	cases := []struct {
		description string
		factory     Factory
		args        []interface{}
		code        int
		reason      string
		msg         string
		data        map[string]interface{}
	}{
		{
			description: "error can be new right with one arg",
			factory: MustNewFactory(http.StatusBadRequest,
				"FormatError",
				"Expected format is '[a-z]*', actual %{format}",
			),
			args: []interface{}{
				"a0",
			},
			code:   http.StatusBadRequest,
			reason: "FormatError",
			msg:    "Expected format is '[a-z]*', actual a0",
			data: map[string]interface{}{
				"format": "a0",
			},
		},
		{
			description: "error can be new with missing arg",
			factory: MustNewFactory(http.StatusBadRequest,
				"FormatError",
				"Expected format is '[a-z]*', actual %{format}",
			),
			args:   []interface{}{},
			code:   http.StatusBadRequest,
			reason: "FormatError",
			msg:    "Expected format is '[a-z]*', actual %!v(MISSING)",
			data:   map[string]interface{}{},
		},
		{
			description: "error can be new with extra arg",
			factory: MustNewFactory(http.StatusBadRequest,
				"FormatError",
				"Expected format is '[a-z]*', actual %{format}",
			),
			args: []interface{}{
				"a0",
				"another",
			},
			code:   http.StatusBadRequest,
			reason: "FormatError",
			msg:    "Expected format is '[a-z]*', actual a0%!(EXTRA string=another)",
			data: map[string]interface{}{
				"format":         "a0",
				"(EXTRA)another": "another",
			},
		},
	}

	for _, c := range cases {
		err := c.factory.New(c.args...)
		e, ok := err.(*Error)
		assert.True(t, ok, c.description)
		assert.Equal(t, c.code, e.Code, c.description, ": code should be equal")
		assert.Equal(t, c.reason, e.Reason, c.description, ": reason should be equal")
		assert.Equal(t, c.msg, e.Message, c.description, ": message should be equal")
		assert.Equal(t, c.data, e.Data, c.description, ": data should be equal")
	}
}
