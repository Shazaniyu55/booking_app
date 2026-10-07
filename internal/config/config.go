package config

import "os"

type Config struct {
	Port     string
	MongoURI string
	MongoDB  string
}

func Load() Config {
	return Config{
		Port:     get("PORT", "8080"),
		MongoURI: get("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  get("MONGO_DB", "booking_app"),
	}
}

func get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
