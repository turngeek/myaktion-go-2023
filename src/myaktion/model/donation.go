package model

type Status string

const (
	TRANSFERRED Status = "TRANSFERRED"
	IN_PROCESS  Status = "IN_PROCESS"
)

type Donation struct {
	Amount           float64
	DonorName        string
	ReceiptRequested bool
	Status           Status
	Account          Account
}
