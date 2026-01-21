package api

import (
	"log"
	"net/http"
	"nofx/config"
	"nofx/crypto"

	"github.com/gin-gonic/gin"
)

// CryptoHandler Encryption API handler
type CryptoHandler struct {
	cryptoService *crypto.CryptoService
}

// NewCryptoHandler Creates encryption handler
func NewCryptoHandler(cryptoService *crypto.CryptoService) *CryptoHandler {
	return &CryptoHandler{
		cryptoService: cryptoService,
	}
}

// ==================== Crypto Config Endpoint ====================

// HandleGetCryptoConfig Get crypto configuration
func (h *CryptoHandler) HandleGetCryptoConfig(c *gin.Context) {
	cfg := config.Get()
	c.JSON(http.StatusOK, gin.H{
		"transport_encryption": cfg.TransportEncryption,
	})
}

// ==================== Public Key Endpoint ====================

// HandleGetPublicKey Get server public key
func (h *CryptoHandler) HandleGetPublicKey(c *gin.Context) {
	cfg := config.Get()
	if !cfg.TransportEncryption {
		c.JSON(http.StatusOK, gin.H{
			"public_key":           "",
			"algorithm":            "",
			"transport_encryption": false,
		})
		return
	}

	publicKey := h.cryptoService.GetPublicKeyPEM()
	c.JSON(http.StatusOK, gin.H{
		"public_key":           publicKey,
		"algorithm":            "RSA-OAEP-2048",
		"transport_encryption": true,
	})
}

// ==================== Encrypted Data Decryption Endpoint ====================

// HandleDecryptSensitiveData Decrypt encrypted data sent from client
func (h *CryptoHandler) HandleDecryptSensitiveData(c *gin.Context) {
	var payload crypto.EncryptedPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Decrypt
	decrypted, err := h.cryptoService.DecryptSensitiveData(&payload)
	if err != nil {
		log.Printf("❌ Decryption failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failed"})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"plaintext": decrypted,
	})
}

// ==================== Audit Log Query Endpoint ====================

// Audit log functionality removed, not needed in current simplified implementation

// ==================== Utility Functions ====================

// isValidPrivateKey Validate private key format
func isValidPrivateKey(key string) bool {
	// EVM private key: 64 hex characters (optional 0x prefix)
	if len(key) == 64 || (len(key) == 66 && key[:2] == "0x") {
		return true
	}
	// TODO: Add validation for other chains
	return false
}

// Handler is the entry point for Vercel Serverless Functions
// This function is required by Vercel to recognize this file as a serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	// For Vercel deployment, we need to route requests to the appropriate handler
	// Since this is part of a larger Gin application, we create a minimal Gin router
	// to handle crypto-related routes
	
	// Create a new Gin context from the HTTP request
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	
	// Note: This is a minimal implementation for Vercel compatibility
	// In a full deployment, you would need to initialize the CryptoHandler instance
	// with a cryptoService dependency
	// For now, this provides a basic handler that Vercel can recognize
	
	// Route to a simple handler that indicates the function is available
	router.Any("/*path", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Crypto API handler is available",
			"path":    c.Request.URL.Path,
		})
	})
	
	// Serve the request
	router.ServeHTTP(w, r)
}
