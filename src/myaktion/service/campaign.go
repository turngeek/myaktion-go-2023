package service

import (
	"log"

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
	log.Printf("Successfully stored new campaign with ID %v in database.", campaign.ID)
	log.Printf("Stored: %v", campaign)
	return nil
}
