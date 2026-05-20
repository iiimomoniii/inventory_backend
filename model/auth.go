package model

import "time"

// ─── Request ───────────────────────────────────────────────

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// ─── Response ──────────────────────────────────────────────

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
}

// ─── Entity ────────────────────────────────────────────────

type RefreshToken struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"userId"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
	RevokedAt *time.Time `json:"revokedAt"`
}
