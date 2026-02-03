package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"api_kino/service/env"
)

type ConfigStruct struct {
	AppName               string
	AppPath               string
	Key                   string
	Port                  string
	FcmKey                string
	GinMode               string
	RedisHost             string
	RedisPort             string
	AuthApiKey            string
	EnableJob             bool
	EnableFreeMem         bool
	DBMaxIdleCon          int
	DBMaxOpenCon          int
	DBMaxLifeTime         int
	DBMaxIdleTime         int
	TsskForwardByFileID   string
	TsskForwardJsonPath   string
	BroadcastRegistration bool
	BroadcastFinish       bool
	SSHEnable             bool
	BaseUrl               string
	BaseUrlHTPDF          string
	BaseBinkarApiUrl      string
	SpersadPhotoPath      string
	GenerateModel         bool
}

type ConfigDB struct {
	Postgres *gorm.DB
	Redis    *redis.Client
}

func Config() ConfigStruct {
	return ConfigStruct{
		AppName:               env.GetEnv("APP_NAME"),
		AppPath:               env.GetEnv("APP_PATH"),
		Key:                   env.GetEnv("APP_KEY"),
		Port:                  env.GetEnv("APP_PORT"),
		FcmKey:                env.GetEnv("FCM_KEY"),
		GinMode:               env.GetEnv("GIN_MODE", gin.ReleaseMode),
		RedisHost:             env.GetEnv("REDIS_HOST"),
		RedisPort:             env.GetEnv("REDIS_PORT"),
		AuthApiKey:            env.GetEnv("AUTH_API_KEY", "e3460062-6a33-11ed-89f2-d336c711220a"),
		EnableJob:             env.GetEnvBool("ENABLE_JOB", true),
		EnableFreeMem:         env.GetEnvBool("ENABLE_FREE_MEM", true),
		DBMaxOpenCon:          env.GetEnvInt("DB_MAX_OPEN_CON", 50),
		DBMaxIdleCon:          env.GetEnvInt("DB_MAX_IDLE_CON", 10),
		DBMaxLifeTime:         env.GetEnvInt("DB_MAX_LIFE_TIME", 50),
		DBMaxIdleTime:         env.GetEnvInt("DB_MAX_IDLE_TIME", 50),
		TsskForwardByFileID:   env.GetEnv("TSSK_FORWARD_BY_FILE_ID", ""),
		TsskForwardJsonPath:   env.GetEnv("TSSK_FORWARD_JSON_PATH", ""),
		BroadcastRegistration: env.GetEnvBool("BROADCAST_REGISTRATION", false),
		BroadcastFinish:       env.GetEnvBool("BROADCAST_FINISH", false),
		SSHEnable:             env.GetEnvBool("SSH_ENABLE", false),
		BaseUrl:               env.GetEnv("BASE_URL"),
		BaseUrlHTPDF:          env.GetEnv("BASE_URL_HTPDF"),
		BaseBinkarApiUrl:      env.GetEnv("BASE_BINKAR_API_URL", "http://127.0.0.1:14021/v1"),
		SpersadPhotoPath:      env.GetEnv("SPERSAD_PHOTO_PATH", "http://172.27.27.187:8765/ad-fileupload-service/%s/noPers/%s/noTam/%s/download"),
		GenerateModel:         env.GetEnvBool("GENERATE_MODEL", false),
	}
}
