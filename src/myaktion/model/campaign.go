package model

type Campaign struct {
	ID                 uint
	Name               string
	OrganizerName      string
	TargetAmount       float64
	DonationMinimum    float64
	Donations          []Donation
	AmountDonatedSoFar float64
	Account            Account
}
