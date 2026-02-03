package env

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strconv"
)

func LoadEnv() {
	if err := godotenv.Load(); err == nil {
		return
	}
	ex, err := os.Executable()
	if err != nil {
		print(err)
	}
	envPath := filepath.Dir(ex) + "/.env"
	if err := godotenv.Load(envPath); err != nil {
		//log.Print("No .env file found", envPath)
	}
}

func GetEnv(key string, def ...string) string {
	LoadEnv()
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func GetEnvInt(name string, def int) int {
	s := GetEnv(name)
	i, err := strconv.ParseInt(s, 10, 0)
	if nil != err {
		return def
	}
	return int(i)
}

func GetEnvBool(key string, def bool) bool {
	s1 := GetEnv(key)
	result := def
	if s1 != "" {
		result, _ = strconv.ParseBool(s1)
	}
	return result
}

func SetEnv(key string, value string) {
	LoadEnv()
	_ = os.Setenv(key, value)
}
