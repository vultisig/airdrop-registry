package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/vultisig/mobile-tss-lib/tss"

	"github.com/sirupsen/logrus"

	cache "github.com/patrickmn/go-cache"
	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
)

// Api is the main handler for the API
type Api struct {
	logger        *logrus.Logger
	cfg           *config.Config
	s             *services.Storage
	router        *gin.Engine
	openseaAPIKey string
	cachedData    *cache.Cache
}

// NewApi creates a new Api instance
func NewApi(cfg *config.Config, s *services.Storage) (*Api, error) {
	if nil == cfg {
		return nil, fmt.Errorf("config is nil")
	}
	if nil == s {
		return nil, fmt.Errorf("storage is nil")
	}
	return &Api{
		cfg:           cfg,
		s:             s,
		router:        gin.Default(),
		logger:        logrus.WithField("module", "api").Logger,
		openseaAPIKey: cfg.OpenSea.APIKey,
		cachedData:    cache.New(5*time.Minute, 10*time.Minute),
	}, nil
}

func (a *Api) setupRouting() {
	a.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace with your allowed origins
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "x-hex-chain-code"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	a.router.Use(gzip.Gzip(gzip.DefaultCompression))
	a.router.Use(ErrorHandler())
	// register api group
	rg := a.router.Group("/api")
	// endpoint for health check
	rg.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Vultisig Airdrop Registry",
		})
	})
	// Derive PublicKey
	rg.POST("/derive-public-key", a.derivePublicKeyHandler)
	// Vaults
	rg.POST("/vault", a.registerVaultHandler)
	rg.DELETE("/vault/:ecdsaPublicKey/:eddsaPublicKey", a.deleteVaultHandler)
	rg.GET("/vault/:ecdsaPublicKey/:eddsaPublicKey", a.getVaultHandler)
	rg.POST("/vault/:ecdsaPublicKey/:eddsaPublicKey/alias", a.updateAliasHandler)
	rg.GET("/vault/shared/:uid", a.getVaultByUIDHandler)
	rg.POST("/vault/join-airdrop", a.joinAirdrop)
	rg.POST("/vault/exit-airdrop", a.exitAirdrop)

	// Coins
	rg.DELETE("/coin/:ecdsaPublicKey/:eddsaPublicKey/:coinID", a.deleteCoin)
	rg.POST("/coin/:ecdsaPublicKey/:eddsaPublicKey", a.addCoin)
	rg.GET("/coin/:ecdsaPublicKey/:eddsaPublicKey", a.getCoin)

	// Vault Share Appearance
	rg.GET("vault/theme/:uid", a.getVaultShareAppearanceHandler)
	rg.POST("vault/theme", a.updateVaultShareAppearanceHandler)

	// leader board
	rg.GET("/leaderboard/vaults", a.getVaultsByRankHandler)

	// NFT-related endpoints
	rg.GET("/nft/price/:collectionID", a.getCollectionMinPriceHandler)
	rg.POST("/nft/avatar", a.setNftAvatarHandler)

}

func (a *Api) Start() error {
	a.setupRouting()
	return a.router.Run(fmt.Sprintf("%s:%d", a.cfg.Server.Host, a.cfg.Server.Port))
}

func (a *Api) derivePublicKeyHandler(c *gin.Context) {
	var req models.DerivePublicKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Error(err)
		c.Error(errInvalidRequest)
		return
	}
	result, err := tss.GetDerivedPubKey(req.PublicKeyECDSA, req.HexChainCode, req.DerivePath, false)
	if err != nil {
		a.logger.Error(err)
		c.Error(errFailedToDerivePublicKey)
		return
	}
	c.JSON(http.StatusOK, gin.H{"public_key": result})
}
