package main

import (
	"fmt"

	"github.com/turngeek/myaktion-go-2023/src/myaktion/model"
)

func main() {
	c := model.Campaign{
		Name:            "Kinder helfen",
		OrganizerName:   "Hans Schmidt",
		TargetAmount:    10000.0,
		DonationMinimum: 10.0,
		Account: model.Account{
			Name:     "Kinder helfen",
			BankName: "Raiffeisen",
			Number:   "1234567890",
		},
	}
	c.Donations = append(c.Donations, model.Donation{
		Amount:           100.0,
		DonorName:        "Hans Muster",
		ReceiptRequested: true,
		Status:           model.IN_PROCESS,
		Account: model.Account{
			Name:     "Hans Muster",
			BankName: "Raiffeisen",
			Number:   "1234567890",
		},
	})
	fmt.Printf("%+v\n", c)
}
