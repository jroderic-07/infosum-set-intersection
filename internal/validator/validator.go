package validator

// Interface defines how keys are validated once read from the data file.
// Implemented by RegexValidator to use defined regex patterns to validate keys.
// UDPRN constructor function creates a RegexValidator object with only the UDPRN regex pattern.
type KeyValidator interface {
	Validate(key string) error
}
