package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	Port             string `mapstructure:"PORT"`
	DriverServiceURL string `mapstructure:"DRIVER_SERVICE_URL"`
	JWTSecret        string `mapstructure:"JWT_SECRET"`
}

/*
Bu katman, her şeyin başladığı yerdir ve uygulama ayağa kalkmadan önce "nasıl" çalışacağını bilmesi gerekir.
Uygulamanın konfigürasyonları bu katmanda yönetilir. Adı üstünde "config".
PORT, DRIVER_SERVICE_URL, JWT_SECRET(Güvenlik Anahtarı) gibi kritik bilgileri ".env" dosyasından okur. Viper kütüphanesi sayesinde bu veriler okunur ve doğrulanır. Bu yapı
olmadan gateway nereye bağlanacağını bilemez.
*/
func LoadConfig() (Config, error) {

	var cfg Config

	// 1. Önce .env dosyasını okumayı dene (Local ortam için)
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		// .env dosyası yoksa bile takılmadan devam eder
	}

	// 2. ortam değişkenlerini bağla
	viper.AutomaticEnv()

	viper.BindEnv("PORT")
	viper.BindEnv("DRIVER_SERVICE_URL")
	viper.BindEnv("JWT_SECRET")

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	if cfg.DriverServiceURL == "" {
		return cfg, errors.New("DRIVER_SERVICE_URL konfigürasyonu yapılmalı.")
	}

	return cfg, nil

}
