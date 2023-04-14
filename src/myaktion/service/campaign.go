package service

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/turngeek/myaktion-go-2023/src/myaktion/db"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/model"
)

var (
	campaignStore map[uint]*model.Campaign
	actCampaignId uint = 1
)

func init() {
	campaignStore = make(map[uint]*model.Campaign)
}

func CreateCampaign(campaign *model.Campaign) error {
	result := db.DB.Create(campaign)
	if result.Error != nil {
		return result.Error
	}
	log.Infof("Successfully stored new campaign with ID %v in database.", campaign.ID)
	log.Tracef("Stored: %v", campaign)
	return nil
}

func GetCampaigns() ([]model.Campaign, error) {
	var campaigns []model.Campaign
	result := db.DB.Preload("Donations").Find(&campaigns)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Tracef("Retrieved: %v", campaigns)
	return campaigns, nil
}

func GetCampaign(id uint) (*model.Campaign, error) {
	campaign := campaignStore[id]
	if campaign == nil {
		return nil, fmt.Errorf("no campaign with ID %d", id)
	}
	log.Tracef("Retrieved: %v", campaign)
	return campaign, nil
}

func UpdateCampaign(id uint, campaign *model.Campaign) (*model.Campaign, error) {
	existingCampaign, err := GetCampaign(id)
	if err != nil {
		return existingCampaign, err
	}
	existingCampaign.Name = campaign.Name
	existingCampaign.OrganizerName = campaign.OrganizerName
	existingCampaign.TargetAmount = campaign.TargetAmount
	existingCampaign.DonationMinimum = campaign.DonationMinimum
	entry := log.WithField("ID", id)
	entry.Info("Successfully updated campaign.")
	entry.Tracef("Updated: %v", campaign)
	return existingCampaign, nil
}

func DeleteCampaign(id uint) (*model.Campaign, error) {
	campaign, err := GetCampaign(id)
	if err != nil {
		return campaign, err
	}
	delete(campaignStore, id)
	entry := log.WithField("ID", id)
	entry.Info("Successfully deleted campaign.")
	entry.Tracef("Deleted: %v", campaign)
	return campaign, nil
}
