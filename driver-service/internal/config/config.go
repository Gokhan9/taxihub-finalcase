package config

import (
	"github.com/spf13/viper"
)

/*
TODO: ".env" dosyasında tutacağımız değerlerin GO tarafında ki karşılıkları. 'mapstructure' taglari ile Viper içerisinde ki anahtarları struct alanlarıyla eşlemek.
*/
type Config struct {
	Port        string `mapstructure:"PORT"`
	MongoURI    string `mapstructure:"MONGO_URI"`
	MongoDBName string `mapstructure:"MONGO_DB_NAME"`
}

func LoadConfig() (Config, error) {

	var cfg Config

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		// .env dosyası olmasa bile devam etmesini sağlarız. (docker)
	}

	// 2. Ortam değişkenlerini(docker) bağla.
	viper.AutomaticEnv()

	viper.BindEnv("PORT")
	viper.BindEnv("MONGO_URI")
	viper.BindEnv("MONGO_DB_NAME")

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
