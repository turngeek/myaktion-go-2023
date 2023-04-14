package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/model"
)

func AddDonation(campaignId uint, donation *model.Donation) error {
	campaign, err := GetCampaign(campaignId)
	if err != nil {
		return err
	}
	campaign.Donations = append(campaign.Donations, *donation)
	entry := log.WithField("ID", campaignId)
	entry.Info("Successfully added new donation to campaign.")
	entry.Tracef("Stored: %v", donation)
	return nil
}
