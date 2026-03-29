package jwt

//实现生成和解析 Token

import (
	"errors"
	"time"

	"geekedu-project/common/config"

	"github.com/golang-jwt/jwt/v4"
)

// 自定义Claims（保存用户ID、用户名和角色）
type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// 生成Token
func GenerateToken(userID int64, username string, role string) (string, error) {
	cfg := config.InitConfig()
	// 设置过期时间
	expireTime := time.Now().Add(time.Duration(cfg.JWTExpire) * time.Second)

	// 构造Claims
	claims := &CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "geekedu",
		},
	}

	// 签名生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// 解析Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	cfg := config.InitConfig()
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// 验证Token并提取Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
