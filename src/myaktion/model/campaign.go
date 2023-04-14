package model

import "gorm.io/gorm"

type Campaign struct {
	gorm.Model
	Name               string     `gorm:"notNull;size:30"`
	OrganizerName      string     `gorm:"notNull"`
	TargetAmount       float64    `gorm:"notNull;check:target_amount >= 10.0"`
	DonationMinimum    float64    `gorm:"notNull;check:donation_minimum >= 1.0"`
	Donations          []Donation `gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AmountDonatedSoFar float64    `gorm:"-"`
	Account            Account    `gorm:"embedded;embeddedPrefix:account_"`
}

func (c *Campaign) AfterFind(tx *gorm.DB) (err error) {
	var sum float64
	result := tx.Model(&Donation{}).Select("ifnull(sum(amount),0)").Where("campaign_id = ?", c.ID).Scan(&sum)
	if result.Error != nil {
		return result.Error
	}
	c.AmountDonatedSoFar = sum
	return nil
}
