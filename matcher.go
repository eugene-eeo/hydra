package main

import "regexp"
import "errors"

var InvalidMatcher error = errors.New("invalid matcher")

type Matcher interface {
	Match(string) bool
}

type OrMatcher []Matcher
type AndMatcher []Matcher
type RegexMatcher struct {
	r *regexp.Regexp
}

func (a AndMatcher) Match(s string) bool {
	for _, x := range a {
		if !x.Match(s) {
			return false
		}
	}
	return true
}

func (a OrMatcher) Match(s string) bool {
	for _, x := range a {
		if x.Match(s) {
			return true
		}
	}
	return false
}

func (r *RegexMatcher) Match(s string) bool {
	return r.r.MatchString(s)
}

func interfaceToMatcher(v interface{}) (Matcher, error) {
	switch v.(type) {
	case string:
		return &RegexMatcher{regexp.MustCompile(v.(string))}, nil
	case map[string]interface{}:
		m := v.(map[string]interface{})
		is_and := true
		var a interface{} = nil
		if _, ok := m["&&"]; ok {
			is_and = true
			a = m["&&"]
		} else if _, ok := m["||"]; ok {
			is_and = false
			a = m["||"]
		} else {
			return nil, InvalidMatcher
		}
		arr, ok := a.([]interface{})
		if !ok {
			return nil, InvalidMatcher
		}
		matchers := make([]Matcher, len(arr))
		for i, x := range arr {
			matcher, err := interfaceToMatcher(x)
			if err != nil {
				return nil, err
			}
			matchers[i] = matcher
		}
		if is_and {
			return AndMatcher(matchers), nil
		} else {
			return OrMatcher(matchers), nil
		}
	}
	return nil, InvalidMatcher
}
