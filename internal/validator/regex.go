package validator

import (
	"fmt"
	"regexp"
)

type RegexValidator struct {
	patterns []*regexp.Regexp
}

func NewRegexValidator(patterns ...string) (*RegexValidator, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("compile pattern %q: %w", pattern, err)
		}
		compiled = append(compiled, re)
	}

	return &RegexValidator{patterns: compiled}, nil
}

func (v *RegexValidator) Validate(key string) error {
	if key == "" {
		return fmt.Errorf("empty key")
	}

	for _, pattern := range v.patterns {
		if pattern.MatchString(key) {
			return nil
		}
	}

	return fmt.Errorf("key %q did not match any accepted format", key)
}
