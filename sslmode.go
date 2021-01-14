package postgres

// SSLMode enums.
type SSLMode string

const (
	SSLModeDisable    SSLMode = "disable"
	SSLModeVerifyCA           = "verify-ca"
	SSLModeVerifyFull         = "verify-full"
)

func (ssl SSLMode) String() string {
	return string(ssl)
}
