package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	errUnknown                 = errors.New("UNKNOWN_ERROR")
	errInvalidRequest          = errors.New("INVALID_REQUEST")
	errVaultAlreadyRegist      = errors.New("VAULT_ALREADY_REGISTERED")
	errFailedToRegisterVault   = errors.New("FAIL_TO_REGISTER_VAULT")
	errVaultNotFound           = errors.New("VAULT_NOT_FOUND")
	errFailedToGetVault        = errors.New("FAIL_TO_GET_VAULT")
	errFailedToDeleteVault     = errors.New("FAIL_TO_DELETE_VAULT")
	errFailedToGetCoin         = errors.New("FAIL_TO_GET_COIN")
	errFailedToJoinRegistry    = errors.New("FAIL_TO_JOIN_REGISTRY")
	errFailedToExitRegistry    = errors.New("FAIL_TO_EXIT_REGISTRY")
	errForbiddenAccess         = errors.New("FORBIDDEN_ACCESS")
	errFailedToGetAddress      = errors.New("FAIL_TO_GET_ADDRESS")
	errAddressNotMatch         = errors.New("ADDRESS_NOT_MATCH")
	errFailedToAddCoin         = errors.New("FAIL_TO_ADD_COIN")
	errFailedToDeleteCoin      = errors.New("FAIL_TO_DELETE_COIN")
	errFailedToDerivePublicKey = errors.New("FAIL_TO_DERIVE_PUBLIC_KEY")
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var statusCode int
			errText := err.Error()
			switch err {
			case errInvalidRequest, errVaultAlreadyRegist:
				statusCode = http.StatusBadRequest
			case errAddressNotMatch:
				statusCode = http.StatusBadRequest
			case errVaultNotFound:
				statusCode = http.StatusNotFound
			case errForbiddenAccess:
				statusCode = http.StatusForbidden
			case errFailedToRegisterVault, errFailedToGetVault, errFailedToDeleteVault,
				errFailedToGetCoin, errFailedToJoinRegistry, errFailedToExitRegistry,
				errFailedToGetAddress, errFailedToAddCoin, errFailedToDeleteCoin,
				errFailedToDerivePublicKey:
				statusCode = http.StatusInternalServerError

			default:
				statusCode = http.StatusInternalServerError
				errText = errUnknown.Error()
			}
			c.JSON(statusCode, gin.H{"error": errText})
			c.Abort()
			return
		}
	}
}