package simplpay

import "go.caringcompany.co/simpliclaim/utils/timeops"

// DatabaseProvider provides methods for performing operations on scheduler.
type DatabaseProvider interface {
	UserStore
	MerchantStore
	TransactionStore
}

var (
	// Now func which provide time functionality.
	Now = timeops.Now
)


// Service sets up the basic service.
type Service struct {
	DB DatabaseProvider
	// UUID uuid.Provider
	// Now  timeops.Now (create provider)
}
