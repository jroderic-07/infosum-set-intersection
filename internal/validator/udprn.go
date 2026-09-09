package validator

const udprnPattern = `^\d{8}$`

func NewUDPRNValidator() (*RegexValidator, error) {
	return NewRegexValidator(udprnPattern)
}
