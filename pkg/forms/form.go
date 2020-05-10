package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailReg = regexp.MustCompile("^[\\d\\w.-]+@[\\d\\w.-]+\\.\\w{2,}$")

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		map[string][]string{},
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "Field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, l int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > l {
		f.Errors.Add(field, fmt.Sprintf("Field can not be longer than %d characters", l))
	}
}

func (f *Form) MinLength(field string, l int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) < l {
		f.Errors.Add(field, fmt.Sprintf("Field can not be shorter than %d characters", l))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		f.Errors.Add(field, "Field is invalid")
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "Field is invalid")
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
