package paymentmethod

import "streetbox.id/model"

// RepoInterface ...
type RepoInterface interface {
	FindByActive() *[]model.ResPaymentMethod
}
