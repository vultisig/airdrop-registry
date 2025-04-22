package vultibot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/models"
)

type VultibotService struct {
	baseAddress string
	apiKey      string
	logger      *logrus.Logger
}

func NewVultibotService(baseAddress, apiKey string) *VultibotService {
	return &VultibotService{
		baseAddress: baseAddress,
		apiKey:      apiKey,
		logger:      logrus.New(),
	}
}

func (v *VultibotService) GetReferrals(eddsaKey string, ecdsaKey string) ([]models.Referral, error) {
	url := fmt.Sprintf("%s/user/referrals?eddsaKey=%s&ecdsaKey=%s&apiKey=%s",
		v.baseAddress,
		eddsaKey,
		ecdsaKey,
		v.apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		v.logger.WithError(err).Error("Failed to fetch referrals from API")
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		v.logger.Errorf("API returned non-200 status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var apiResponse models.ReferralsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		v.logger.WithError(err).Error("Failed to fetch referrals from API")
		return nil, err
	}
	return apiResponse.Items, nil
}

func (v *VultibotService) GetAllAchievements(achievementsRequest models.AchievementsRequest) ([]models.AchievementsResponse, error) {
	url := fmt.Sprintf("%s/achievements/list",
		v.baseAddress,
	)

	reqBody, err := json.Marshal(achievementsRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		v.logger.WithError(err).Error("Failed to fetch from API")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching Achivements : %s", resp.Status)
	}

	var response []models.AchievementsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		v.logger.WithError(err).Error("Failed to decode achievements from API")
		return nil, fmt.Errorf("error decoding achievements response: %w", err)
	}

	return response, nil
}
