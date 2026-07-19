package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
	RedisUri     string `mapstructure:"REDIS_URL"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     string `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`
	FrontendURL  string `mapstructure:"FRONTEND_URL"`

	ResendAPIKey string `mapstructure:"RESEND_API_KEY"`
	EmailFrom    string `mapstructure:"EMAIL_FROM"`

	// SIS (Shifd Labs Identity Service) — Phase 2 JWKS-based auth.
	SISJWKSURL          string        `mapstructure:"SIS_JWKS_URL"`
	SISIssuer           string        `mapstructure:"SIS_ISSUER"`
	JWKSRefreshInterval time.Duration `mapstructure:"JWKS_REFRESH_INTERVAL"`

	// S3 storage — the backend now holds the credentials and issues short-lived
	// presigned URLs to the browser (AUDIT SEC-01). Any key that was previously
	// shipped in the frontend bundle must be rotated.
	S3Region          string `mapstructure:"S3_REGION"`
	S3Bucket          string `mapstructure:"S3_BUCKET"`
	S3AccessKeyID     string `mapstructure:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey string `mapstructure:"S3_SECRET_ACCESS_KEY"`

	// CORSAllowedOrigins is a comma-separated allowlist of browser origins
	// (AUDIT SEC-10). Empty falls back to the built-in local-dev defaults.
	CORSAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`

	// DBSSLMode selects the Postgres TLS mode (AUDIT SEC-10):
	// disable|require|verify-ca|verify-full. Empty defaults to "disable" to
	// preserve existing local behaviour; set to "require" (or stricter) in prod.
	DBSSLMode string `mapstructure:"POSTGRES_SSLMODE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
