package model

// ─── Request ───────────────────────────────────────────────

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ─── Response ──────────────────────────────────────────────

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"` // seconds
}
