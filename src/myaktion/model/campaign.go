package model

type Campaign struct {
	Name               string     
	OrganizerName      string     
	TargetAmount       float64    
	DonationMinimum    float64    
	Donations          []Donation 
	AmountDonatedSoFar float64    
	Account            Account    
}
