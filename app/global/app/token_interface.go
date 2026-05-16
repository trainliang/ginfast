package app

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenServiceInterface Token服务接口
type TokenServiceInterface interface {
	// GenerateToken 生成JWT令牌
	GenerateToken(user *ClaimsUser) (string, error)

	// ParseToken 解析JWT令牌
	ParseToken(tokenString string) (*Claims, error)

	// ValidateToken 验证JWT令牌
	ValidateToken(tokenString string) (*Claims, error)

	// GenerateTokenWithCache 生成JWT令牌并存储到缓存
	GenerateTokenWithCache(user *ClaimsUser) (string, error)

	// ValidateTokenWithCache 验证JWT令牌（带缓存检查）
	ValidateTokenWithCache(tokenString string) (*Claims, error)

	// RevokeToken 撤销Token（从缓存中移除）
	RevokeTokenWithCache(tokenString string) error

	// GenerateRefreshToken 生成Refresh Token
	GenerateRefreshToken(userID uint) (string, error)

	// GenerateRefreshTokenForUser 生成包含租户上下文的Refresh Token
	GenerateRefreshTokenForUser(user *ClaimsUser) (string, error)

	// ParseRefreshToken 解析Refresh Token
	ParseRefreshToken(tokenString string) (*RefreshTokenClaims, error)

	// ValidateRefreshToken 验证Refresh Token
	ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error)

	// RevokeRefreshToken 撤销Refresh Token
	RevokeRefreshToken(userID uint) error

	// RefreshAccessToken 使用Refresh Token刷新Access Token并记录在缓存中
	RefreshAccessTokenWithCache(refreshTokenString string, user *ClaimsUser) (string, error)

	// RotateRefreshToken 轮换Refresh Token（撤销旧的，生成新的，保持相同的剩余过期时间）
	RotateRefreshToken(oldRefreshToken string) (string, error)
}

// ClaimsUser 用户声明信息
type ClaimsUser struct {
	UserID           uint   `json:"userId"`                     // 用户ID
	Username         string `json:"username"`                   // 用户名
	TenantID         uint   `json:"tenantId,omitempty"`         // 租户ID
	TenantCode       string `json:"tenantCode,omitempty"`       // 租户编码
	Mode             string `json:"mode,omitempty"`             // 租户上下文模式
	AuthSource       string `json:"authSource,omitempty"`       // 认证来源
	IsPlatformAdmin  bool   `json:"isPlatformAdmin,omitempty"`  // 是否平台管理员
	ActorUserID      uint   `json:"actorUserId,omitempty"`      // 真实操作者ID
	ExternalClientID string `json:"externalClientId,omitempty"` // 外部客户端ID
}

// Claims JWT声明结构
type Claims struct {
	ClaimsUser
	jwt.RegisteredClaims
}

// RefreshTokenClaims Refresh Token声明结构
type RefreshTokenClaims struct {
	ClaimsUser
	jwt.RegisteredClaims
}

// RefreshTokenInfo Refresh Token信息
type RefreshTokenInfo struct {
	UserID    uint      `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// TokenInfo Token信息
type TokenInfo struct {
	UserID    uint      `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}
