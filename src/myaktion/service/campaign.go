package service

import (
	"fmt"

	log "github.com/sirupsen/logrus"

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
	campaign.ID = actCampaignId
	campaignStore[actCampaignId] = campaign
	actCampaignId += 1
	log.Infof("Successfully stored new campaign with ID %v in database.", campaign.ID)
	log.Tracef("Stored: %v", campaign)
	return nil
}

func GetCampaigns() ([]model.Campaign, error) {
	var campaigns []model.Campaign
	for _, campaign := range campaignStore {
		// Note: *campaign does only create a shallow copy of the campaign, means the donation slice is still the same as for campaignStore
		// Ok for this use-case as we're just serializing after usage
		campaigns = append(campaigns, *campaign)
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
