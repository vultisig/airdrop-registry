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
	errFailedToUpdateVault     = errors.New("FAIL_TO_UPDATE_VAULT")
	errFailedToExitRegistry    = errors.New("FAIL_TO_EXIT_REGISTRY")
	errForbiddenAccess         = errors.New("FORBIDDEN_ACCESS")
	errFailedToGetAddress      = errors.New("FAIL_TO_GET_ADDRESS")
	errFailedToFetchFromBotApi = errors.New("FAIL_TO_FETCH_FROM_BOT_API")
	errAddressNotMatch         = errors.New("ADDRESS_NOT_MATCH")
	errFailedToAddCoin         = errors.New("FAIL_TO_ADD_COIN")
	errFailedToDeleteCoin      = errors.New("FAIL_TO_DELETE_COIN")
	errFailedToDerivePublicKey = errors.New("FAIL_TO_DERIVE_PUBLIC_KEY")
	errFailedToGetTheme        = errors.New("FAIL_TO_GET_THEME")
	errFailedToSetTheme        = errors.New("FAIL_TO_SET_THEME")
	errLogoTooLarge            = errors.New("LOGO_TOO_LARGE")
	errFailedToGetCollection   = errors.New("FAIL_TO_GET_COLLECTION")
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var statusCode int
			errText := err.Error()
			switch {
			case errors.Is(err, errInvalidRequest),
				errors.Is(err, errVaultAlreadyRegist),
				errors.Is(err, errLogoTooLarge):
				statusCode = http.StatusBadRequest
			case errors.Is(err, errAddressNotMatch):
				statusCode = http.StatusBadRequest
			case errors.Is(err, errVaultNotFound):
				statusCode = http.StatusNotFound
			case errors.Is(err, errForbiddenAccess):
				statusCode = http.StatusForbidden
			case errors.Is(err, errFailedToRegisterVault),
				errors.Is(err, errFailedToGetVault),
				errors.Is(err, errFailedToDeleteVault),
				errors.Is(err, errFailedToGetCoin),
				errors.Is(err, errFailedToJoinRegistry),
				errors.Is(err, errFailedToUpdateVault),
				errors.Is(err, errFailedToExitRegistry),
				errors.Is(err, errFailedToGetAddress),
				errors.Is(err, errFailedToAddCoin),
				errors.Is(err, errFailedToDeleteCoin),
				errors.Is(err, errFailedToDerivePublicKey),
				errors.Is(err, errFailedToSetTheme),
				errors.Is(err, errFailedToFetchFromBotApi),
				errors.Is(err, errFailedToGetTheme),
				errors.Is(err, errFailedToGetCollection):
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
