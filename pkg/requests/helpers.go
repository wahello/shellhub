package requests

import (
	e "github.com/shellhub-io/shellhub/pkg/errors"
	"github.com/shellhub-io/shellhub/pkg/models"
)

func HasBillingInstance(ns *models.Namespace) bool {
	if ns == nil || ns.Billing == nil {
		return false
	}

	return true
}

func HandleStatusResponse(status int) error {
	switch status {
	case 200:
		return nil
	case 402:
		return nil
	case 400:
		return nil
	case 423:
		return e.ErrLocked
	default:
		return e.ErrReport
	}
}
