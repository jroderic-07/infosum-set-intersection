package validator

import "testing"

func TestRegexValidator_Validate(t *testing.T) {
	validator, err := NewRegexValidator(`^\d{8}$`, `^[A-Z]{2}\d{6}$`)
	if err != nil {
		t.Fatalf("NewRegexValidator returned error: %v", err)
	}

	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{name: "valid udprn", key: "08034283", wantErr: false},
		{name: "valid alternate pattern", key: "AB123456", wantErr: false},
		{name: "empty key", key: "", wantErr: true},
		{name: "header value", key: "udprn", wantErr: true},
		{name: "too short", key: "123", wantErr: true},
		{name: "non numeric", key: "abc12345", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate(%q) error = %v, wantErr %v", tt.key, err, tt.wantErr)
			}
		})
	}
}

func TestNewUDPRNValidator(t *testing.T) {
	validator, err := NewUDPRNValidator()
	if err != nil {
		t.Fatalf("NewUDPRNValidator returned error: %v", err)
	}

	if err := validator.Validate("71842328"); err != nil {
		t.Errorf("Validate(71842328) returned error: %v", err)
	}

	if err := validator.Validate(""); err == nil {
		t.Error("Validate(\"\") should return an error")
	}
}
