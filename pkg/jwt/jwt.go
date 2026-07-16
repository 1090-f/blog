package jwt

import (
	"errors"
	"strings"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// Claims JWT 令牌载荷，包含用户 ID 和角色信息。
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwtv5.RegisteredClaims
}

// JWT 包相关的错误值。
var (
	ErrEmptySecret   = errors.New("jwt secret must not be empty")         // 密钥不能为空
	ErrInvalidExpire = errors.New("jwt expire must be greater than zero") // 过期时间必须大于零
)

// GenerateToken 使用 HS256 算法生成携带用户身份信息的 JWT 令牌。
func GenerateToken(secret string, userID uint, role string, expireSeconds int) (string, error) {
	// 密钥不能为空
	if strings.TrimSpace(secret) == "" {
		return "", ErrEmptySecret
	}
	// 过期时间必须大于零
	if expireSeconds <= 0 {
		return "", ErrInvalidExpire
	}

	// 构建载荷：自定义字段 + 标准字段（过期时间、签发时间）
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	}

	// 使用 HMAC-SHA256 签名
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并校验 JWT 令牌，返回载荷信息；签名无效或已过期时返回错误。
func ParseToken(tokenString string, secret string) (*Claims, error) {
	// 密钥不能为空
	if strings.TrimSpace(secret) == "" {
		return nil, ErrEmptySecret
	}

	// 解析令牌并验证签名，keyFunc 返回密钥用于验签
	token, err := jwtv5.ParseWithClaims(tokenString, &Claims{}, func(token *jwtv5.Token) (any, error) {
		// 强制校验签名算法必须是 HMAC，防止算法混淆攻击
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, jwtv5.ErrTokenSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	// 提取载荷并确认令牌有效
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwtv5.ErrTokenInvalidClaims
	}

	return claims, nil
}
