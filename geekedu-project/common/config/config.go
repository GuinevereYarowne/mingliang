package config

//全局配置管理，读取环境变量，提供统一访问接口

import (
	"os"
	"strconv"
)

// 全局配置结构体
type Config struct {
	// MySQL配置
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBDatabase string

	// OSS配置
	OSSAccessKey string
	OSSSecretKey string
	OSSEndpoint  string
	OSSBucket    string

	// JWT配置
	JWTSecret string
	JWTExpire int64 // 过期时间（秒）
}

// 初始化配置（从环境变量读取，新手友好）
func InitConfig() *Config {
	// 读取环境变量，设置默认值
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	jwtExpire, _ := strconv.ParseInt(getEnv("JWT_EXPIRE", "86400"), 10, 64)

	return &Config{
		// MySQL
		DBHost:     getEnv("DB_HOST", "mysql"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "root"),
		DBDatabase: getEnv("DB_DATABASE", "geekedu"),

		// OSS（核心：从环境变量读取，不硬编码）
		OSSAccessKey: getEnv("OSS_ACCESS_KEY", ""),
		OSSSecretKey: getEnv("OSS_SECRET_KEY", ""),
		OSSEndpoint:  getEnv("OSS_ENDPOINT", "oss-cn-hangzhou.aliyuncs.com"),
		OSSBucket:    getEnv("OSS_BUCKET", ""),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "geekedu-secret-123"), // 新手可设默认，生产改环境变量
		JWTExpire: jwtExpire,
	}
}

// 辅助函数：读取环境变量，无则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
