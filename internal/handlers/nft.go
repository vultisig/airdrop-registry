package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vultisig/airdrop-registry/internal/common"
)

type SetNftProfileRequest struct {
	Uid            string `json:"uid" binding:"required"`
	PublicKeyECDSA string `json:"public_key_ecdsa" binding:"required"`
	PublicKeyEDDSA string `json:"public_key_eddsa" binding:"required"`
	HexChainCode   string `json:"hex_chain_code" binding:"required"`
	CollectionID   string `json:"collection_id" binding:"required"`
	ItemID         int64  `json:"item_id,string" binding:"required"`
	Url            string `json:"url" binding:"required"`
}

type OpenSeaNFTResponse struct {
	NFT struct {
		Identifier      string `json:"identifier"`
		Collection      string `json:"collection"`
		Contract        string `json:"contract"`
		ImageURL        string `json:"image_url"`
		DisplayImageURL string `json:"display_image_url"`
		MetadataURL     string `json:"metadata_url"`
		Owners          []struct {
			Address  string `json:"address"`
			Quantity int    `json:"quantity"`
		} `json:"owners"`
	} `json:"nft"`
}

func (a *Api) setNftAvatarHandler(c *gin.Context) {
	var vault SetNftProfileRequest
	if err := c.ShouldBindJSON(&vault); err != nil {
		a.logger.Error(err)
		_ = c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	v, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA)
	if err != nil {
		a.logger.Errorf("fail to get vault,err: %v", err)
		_ = c.Error(errFailedToGetVault)
		return
	}
	if v == nil {
		_ = c.Error(errVaultNotFound)
		return
	}
	if v.HexChainCode == vault.HexChainCode && v.Uid == vault.Uid {
		//setProfile(vault.Uid, vault.Url)
		//check if user owns the nft
		var nftOwnerResponse OpenSeaNFTResponse
		key := fmt.Sprintf("%s-%d", vault.CollectionID, vault.ItemID)
		//check cache
		if cachedData, ok := a.cachedData.Get(key); ok {
			if _, ok := cachedData.(OpenSeaNFTResponse); ok {
				nftOwnerResponse = cachedData.(OpenSeaNFTResponse)
			}
		}

		if nftOwnerResponse.NFT.Collection == "" {
			//fetch from opensea
			url := fmt.Sprintf("https://api.opensea.io/api/v2/chain/ethereum/contract/%s/nfts/%d", vault.CollectionID, vault.ItemID)
			// add x-api-key header
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				a.logger.Errorf("failed to create request: %v", err)
				_ = c.Error(errFailedToGetCollection)
				return
			}
			req.Header.Add("x-api-key", a.cfg.OpenSea.APIKey)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				a.logger.Errorf("failed to get collection: %v", err)
				_ = c.Error(errFailedToGetCollection)
				return
			}

			defer a.closer(resp.Body)
			if err := json.NewDecoder(resp.Body).Decode(&nftOwnerResponse); err != nil {
				a.logger.Errorf("failed to decode response: %v", err)
				_ = c.Error(errFailedToGetCollection)
				return
			}
			if err := a.cachedData.Add(key, nftOwnerResponse, time.Minute); err != nil {
				a.logger.Errorf("fail to add collection to cache: %s", err)
			}
		}

		owned := false
		ethAddress, err := v.GetAddress(common.Ethereum)
		if err != nil {
			a.logger.Errorf("fail to get address: %v", err)
			_ = c.Error(errFailedToGetAddress)
			return
		}
		if nftOwnerResponse.NFT.Owners != nil {
			for _, owner := range nftOwnerResponse.NFT.Owners {
				if strings.EqualFold(owner.Address, ethAddress) {
					owned = true
					break
				}
			}
		}
		if !owned {
			_ = c.Error(errForbiddenAccess)
			return
		}
		v.AvatarCollectionID = vault.CollectionID
		v.AvatarItemID = vault.ItemID
		v.AvatarURL = vault.Url
		if err := a.s.UpdateVaultAvatar(v); err != nil {
			a.logger.Errorf("fail to update vault avatar: %v", err)
			_ = c.Error(err)
			return
		}
	} else {
		_ = c.Error(errForbiddenAccess)
		return
	}
	c.Status(http.StatusOK)
}

type OpenSeaBestCollectionResponse struct {
	Listings []struct {
		Price struct {
			Current struct {
				Currency string `json:"currency"`
				Decimals int    `json:"decimals"`
				Value    string `json:"value"`
			} `json:"current"`
		} `json:"price"`
	} `json:"listings"`
}

func (a *Api) getCollectionMinPriceHandler(c *gin.Context) {
	collectionID := c.Param("collectionID")
	if collectionID == "" {
		_ = c.Error(errInvalidRequest)
		return
	}
	if !strings.EqualFold(collectionID, "0xa98b29a8f5a247802149c268ecf860b8308b7291") {
		_ = c.Error(errAddressNotMatch)
		return
	}
	collectionSlug := "thorguards"
	//check cache first
	if cached, ok := a.cachedData.Get(collectionSlug); ok {
		if price, ok := cached.(OpenSeaBestCollectionResponse); ok {
			c.JSON(http.StatusOK, gin.H{"minPrice": price.Listings[0].Price.Current})
			return
		}
	}
	url := fmt.Sprintf("https://api.opensea.io/api/v2/listings/collection/%s/best", collectionSlug)
	// add x-api-key header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		a.logger.Errorf("failed to create request: %v", err)
		_ = c.Error(errFailedToGetCollection)
		return
	}
	a.logger.Infof("api key: %s", a.cfg.OpenSea.APIKey)
	req.Header.Add("x-api-key", a.cfg.OpenSea.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.logger.Errorf("failed to get collection: %v", err)
		_ = c.Error(errFailedToGetCollection)
		return
	}
	defer a.closer(resp.Body)
	var openseaResp OpenSeaBestCollectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&openseaResp); err != nil {
		a.logger.Errorf("failed to decode response: %v", err)
		_ = c.Error(errFailedToGetCollection)
		return
	}
	if len(openseaResp.Listings) == 0 {
		_ = c.Error(errFailedToGetCollection)
		return
	}
	//add to cache
	if err := a.cachedData.Add(collectionSlug, openseaResp.Listings[0].Price, time.Minute); err != nil {
		a.logger.Errorf("fail to add collection to cache: %s", err)
	}
	c.JSON(http.StatusOK, gin.H{"minPrice": openseaResp.Listings[0].Price.Current})
}
