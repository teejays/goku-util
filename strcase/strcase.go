package strcase

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/teejays/goku-util/panics"
)

var acronyms = map[string]bool{
	"ID":   true,
	"UUID": true,
	"UI":   true,
	"DAL":  true,
	"USA":  true,
	"HTTP": true,
	"API":  true,
	"JWT":  true,
}

func init() {
	// Note: Not sure if this even works. Our own Acronyms work better
	for k, v := range acronyms {
		if !v {
			continue
		}
		strcase.ConfigureAcronym(k, strings.ToLower(k))
	}
}

func IsEqual(a, b string) bool {
	return strings.EqualFold(a, b)
}
func IsAcronym(s string) bool {
	return acronyms[strings.ToUpper(s)]
}

func HasAcronym(s string) bool {
	for acr, is := range acronyms {
		if is && strings.Contains(s, acr) {
			return true
		}
	}
	return false
}

// UpperizeAcronym - only upperize if the prefix or suffix is an Acronym
func UpperizeAcronym(s string) string {
	for acr, is := range acronyms {
		if !is {
			continue
		}
		// Can't match if string is smaller than the acronym
		if len(s) < len(acr) {
			continue
		}
		// First n characters
		if strings.EqualFold(s[:len(acr)], acr) {
			s = strings.ToUpper(s[:len(acr)]) + s[len(acr):]
		}
		// Last n characters
		if strings.EqualFold(s[len(s)-len(acr):], acr) {
			s = s[:len(s)-len(acr)] + strings.ToUpper(s[len(s)-len(acr):])
		}
		// Last n+1 characters == plural acronym (i.e. xyzs / XYZs)
		if len(s) >= len(acr)+1 &&
			strings.EqualFold(s[len(s)-len(acr)-1:], acr+"s") {
			s = s[:len(s)-len(acr)-1] + strings.ToUpper(s[len(s)-len(acr)-1:len(s)-1]) + "s"
		}
	}
	return s
}

func ToCamel(s string) string {
	if IsAcronym(s) {
		return strings.ToUpper(s)
	}
	s = strcase.ToCamel(s)
	s = UpperizeAcronym(s)
	return s
}

// PartsSep seprates major parts in the string vs. a single `_` which separates words. It is three underscores.
const PartsSep = "___"

// to___split_camel -> To<sep>SplitCamel eg. To.SplitCamel
func ToSplitCamel(s string, sep string) string {
	if IsAcronym(s) {
		return strings.ToUpper(s)
	}
	// split by key
	parts := strings.Split(s, PartsSep)
	for i := range parts {
		parts[i] = strcase.ToCamel(parts[i])
	}
	r := strings.Join(parts, sep)
	return UpperizeAcronym(r)
}

// to___snaked_camel -> To_SnakedCamel
func ToSnakedCamel(s string) string {
	return ToSplitCamel(s, "_")
}

func ToLowerCamel(s string) string {
	if IsAcronym(s) {
		return strings.ToLower(s)
	}
	s = strcase.ToLowerCamel(s)
	s = UpperizeAcronym(s)
	return s
}

func ToSnake(s string) string {
	for acronym, active := range acronyms {
		if !active {
			continue
		}
		// If there is an uppercase acronym XYZ in here, make it title case since strcase thinks x, y, and z in XYX are seperate words
		replace := strings.ToUpper(acronym)
		with := strings.Title(strings.ToLower(acronym))
		s = strings.ReplaceAll(s, replace, with)
	}
	// if there are three `_`, we should maintain them as two `_` in snake case
	parts := strings.Split(s, PartsSep)
	for i := range parts {
		parts[i] = strcase.ToSnake(parts[i])
	}
	return strings.Join(parts, "__")
}

func ToKebab(s string) string {
	return strcase.ToKebab(s)
}

var pluralOverrides = map[string]string{
	"sheep":   "sheep",
	"fish":    "fish",
	"address": "addresses",
	"process": "processes",
}

func Pluralize(s string) string {
	if s == "" {
		return s
	}
	if ans, exists := pluralOverrides[strings.ToLower(s)]; exists {
		return ans
	}
	// If ends with `s`, don't know what to do
	if s[len(s)-1] == 's' {
		panics.P("library strcase: couldn't determine plural form of '%s'", s)
	}
	// If ends with `y`, change `y` to `ies` e.g. company -> companies
	if s[len(s)-1] == 'y' {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}

var singularOverrides = map[string]string{
	"sheep":     "sheep",
	"fish":      "fish",
	"addresses": "address",
	"processes": "process",
}

func Singularize(s string) string {
	if s == "" {
		return s
	}
	if ans, exists := singularOverrides[strings.ToLower(s)]; exists {
		return ans
	}
	if s[len(s)-1] == 's' {
		return s[:len(s)-1]
	}
	panics.P("library strcase: couldn't determine singular form of '%s'", s)
	return s
}
