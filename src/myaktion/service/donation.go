package service

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/client"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/client/banktransfer"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/db"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/model"
)

func AddDonation(campaignId uint, donation *model.Donation) error {
	campaign, err := GetCampaign(campaignId)
	if err != nil {
		return err
	}
	donation.CampaignID = campaign.ID
	result := db.DB.Create(donation)
	if result.Error != nil {
		return result.Error
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := client.GetBankTransferConnection(ctx)
	if err != nil {
		log.Errorf("Failed to connect to bank transfer service: %v", err)
		deleteDonation(donation)
		return err
	}
	defer conn.Close()
	banktransferClient := banktransfer.NewBankTransferClient(conn)
	_, err = banktransferClient.TransferMoney(ctx, &banktransfer.Transaction{
		DonationId:  int32(donation.ID),
		Amount:      float32(donation.Amount),
		Reference:   "Donation",
		FromAccount: convertAccount(&donation.Account),
		ToAccount:   convertAccount(&campaign.Account),
	})
	if err != nil {
		log.Errorf("error calling the banktransfer service: %v", err)
		deleteDonation(donation)
		return err
	}
	entry := log.WithField("ID", campaignId)
	entry.Info("Successfully added new donation to campaign in database.")
	entry.Tracef("Stored: %v", donation)
	return nil
}

func MarkDonation(id uint) error {
	entry := log.WithField("donationId", id)
	donation := new(model.Donation)
	result := db.DB.First(donation, id)
	if result.Error != nil {
		entry.WithError(result.Error).Error("Error retrieving donation")
		return result.Error
	}
	entry = entry.WithField("donation", donation)
	entry.Trace("Retrieved donation")
	donation.Status = model.TRANSFERRED
	result = db.DB.Save(donation)
	if result.Error != nil {
		entry.WithError(result.Error).Error("Can't update donation")
		return result.Error
	}
	entry.Info("Successfully updated status of donation")
	return nil
}

func convertAccount(account *model.Account) *banktransfer.Account {
	return &banktransfer.Account{
		Name:     account.Name,
		BankName: account.BankName,
		Number:   account.Number,
	}
}

func deleteDonation(donation *model.Donation) error {
	entry := log.WithField("donationID", donation.ID)
	entry.Info("Trying to delete donation to make state consistent.")
	result := db.DB.Delete(donation)
	if result.Error != nil {
		// Note: configure logger to raise an alarm to compensate inconsitent state
		entry.WithField("alarm", true).Error("")
		return result.Error
	}
	entry.Info("Successfully deleted donation to make state consistent.")
	return nil
}
