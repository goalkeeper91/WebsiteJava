package domain

// AdminCustomerRow bündelt einen Nutzer mit seiner Subscription (inkl. Tier)
// für die Admin-Kundenliste - Subscription ist nil, falls der Nutzer noch nie
// GetSubscription ausgelöst hat (lazy Free-Tier-Anlage, siehe SubscriptionService),
// in dem Fall ist er implizit ein aktiver Free-Kunde.
type AdminCustomerRow struct {
	User         *User
	Subscription *UserSubscription
}
