package api

type Status struct {
	Status string
	Reason string
	Result interface{}
}

const (
	STATUS_SUCCESS = "STATUS_SUCCESS"
	STATUS_FAILURE = "STATUS_FAILURE"

	REASON_EMPTY             = "REASON_EMPTY"
	REASON_UNKNOWN           = "REASON_UNKNOWN"
	REASON_PARAMETER_MISSING = "REASON_PARAMETER_MISSING"

	AUTH_FAILURE = "AUTH_FAILURE"
)

type Credential struct {
	Token string
}
