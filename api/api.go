package api

import (
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
)

type Status struct {
	Status string
	Reason string
	Result interface{}
}

const (
	// general
	STATUS_SUCCESS = "STATUS_SUCCESS"
	STATUS_FAILURE = "STATUS_FAILURE"

	REASON_EMPTY             = "REASON_EMPTY"
	REASON_UNKNOWN           = "REASON_UNKNOWN"
	REASON_PARAMETER_MISSING = "REASON_PARAMETER_MISSING"
	REASON_NETWORK_UNSTABLE  = "REASON_NETWORK_UNSTABLE"

	// authentication
	AUTH_FAILURE = "AUTH_FAILURE"

	// applications
	APP_NAME_INVALID     = "APP_NAME_INVALID"
	APP_NON_EXIST        = "APP_NON_EXIST"
	APP_ALREADY_EXISTING = "APP_ALREADY_EXISTING"
)

type Credential struct {
	Token string
}

func EnsureRequest(handler func(echo.Context, interface{}) error, request interface{}) func(echo.Context) error {

	return func(c echo.Context) error {

		if readCloser, err := c.Request().GetBody(); err != nil {
			return c.JSON(http.StatusOK, Status{
				Status: STATUS_FAILURE,
				Reason: REASON_NETWORK_UNSTABLE,
			})
		} else if err := json.NewDecoder(readCloser).Decode(&request); err != nil {

			return c.JSON(http.StatusOK, Status{
				Status: STATUS_FAILURE,
				Reason: REASON_PARAMETER_MISSING,
			})

		} else {

			return handler(c, request)

		}

	}
}
