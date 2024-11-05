package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	DbUrl                 string `mapstructure:"DATABASE_CONFIG"`
	Port                  string `mapstructure:"PORT"`
	HashCost              int    `mapstructure:"HASH_COST"`
	Issuer                string `mapstructure:"ISSUER"`
	ExpDurationHour       int    `mapstructure:"EXP_HOUR"`
	JwtSecretKey          string `mapstructure:"JWT_SECRET_KEY"`
	FirebaseDbUrl         string `mapstructure:"FIREBASE_DB_URL"`
	FirebaseProjectId     string `mapstructure:"FIREBASE_PROJECT_ID"`
	FirebaseStorageBucket string `mapstructure:"FIREBASE_STORAGE_BUCKET"`
}

func ConfigInit() (Config, error) {
	env := Config{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return env, nil
}
