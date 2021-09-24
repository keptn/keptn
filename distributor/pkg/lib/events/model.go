package events

// AdditionalSubscriptionData is the data the distributor
// will add as temporary data to the keptn events forwarded
// to the keptn integration
type AdditionalSubscriptionData struct {
	SubscriptionID string `json:"subscriptionID"`
}
