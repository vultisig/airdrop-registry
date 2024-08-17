package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vultisig/mobile-tss-lib/tss"

	"github.com/sirupsen/logrus"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
)

// Api is the main handler for the API
type Api struct {
	logger *logrus.Logger
	cfg    *config.Config
	s      *services.Storage
	router *gin.Engine
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
		cfg:    cfg,
		s:      s,
		router: gin.Default(),
		logger: logrus.WithField("module", "api").Logger,
	}, nil
}

func (a *Api) setupRouting() {
	a.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace with your allowed origins
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	a.router.GET("/webapp", func(c *gin.Context) {
		//server index.html file in demo folder
		c.File("web/dist/index.html")
	})
	a.router.Static("/assets", "web/dist/assets")
	a.router.Static("/fonts", "web/dist/fonts")
	a.router.Static("/images", "web/dist/images")
	a.router.Static("/chains", "web/dist/chains")
	a.router.StaticFile("/wallet-core.wasm", "web/dist/wallet-core.wasm")
	a.router.StaticFile("/favicon.ico", "web/dist/favicon.ico")
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
	rg.DELETE("/vault", a.deleteVaultHandler)
	rg.GET("/vault/:ecdsaPublicKey/:eddsaPublicKey", a.getVaultHandler)
	rg.POST("/vault/join-airdrop", a.joinAirdrop)
	rg.POST("/vault/exit-airdrop", a.exitAirdrop)

	// Coins
	rg.DELETE("/coin/:ecdsaPublicKey/:eddsaPublicKey/:coinID", a.deleteCoin)
	rg.POST("/coin/:ecdsaPublicKey/:eddsaPublicKey", a.addCoin)
	rg.GET("/coin/:ecdsaPublicKey/:eddsaPublicKey", a.getCoin)

}

func (a *Api) Start() error {
	a.setupRouting()
	return a.router.Run(fmt.Sprintf("%s:%d", a.cfg.Server.Host, a.cfg.Server.Port))
}

func (a *Api) derivePublicKeyHandler(c *gin.Context) {
	var req models.DerivePublicKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := tss.GetDerivedPubKey(req.PublicKeyECDSA, req.HexChainCode, req.DerivePath, false)
	if err != nil {
		a.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to derive public key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"public_key": result})
}
