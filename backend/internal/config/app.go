package config

// AppConfig 根配置
type AppConfig struct {
	App        AppConfigApp     `mapstructure:"app" validate:"required"`
	Server     ServerConfig     `mapstructure:"server" validate:"required"`
	Log        LogConfig        `mapstructure:"log" validate:"required"`
	DataSource DataSourceConfig `mapstructure:"datasource" validate:"required"`
	JWT        JWTConfig        `mapstructure:"jwt" validate:"required"`
	Storage    StorageConfig    `mapstructure:"storage" validate:"required"`
}

// AppConfigApp 应用基本信息
type AppConfigApp struct {
	Name string `mapstructure:"name" validate:"required"`
	Env  string `mapstructure:"env" validate:"required,oneof=dev test prod"`
}

// ServerConfig HTTP 服务配置
type ServerConfig struct {
	Addr     string `mapstructure:"addr" validate:"required"`
	Schema   string `mapstructure:"schema" validate:"required,oneof=http https"`
	HostName string `mapstructure:"host_name" validate:"required"`
	Mode     string `mapstructure:"mode" validate:"required,oneof=debug release test"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `mapstructure:"level" validate:"required,oneof=debug info warn error"`
}

// DataSourceConfig 数据源配置
type DataSourceConfig struct {
	Database DatabaseConfig `mapstructure:"database" validate:"required"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver string `mapstructure:"driver" validate:"required,oneof=sqlite mysql postgres"`
	DSN    string `mapstructure:"dsn" validate:"required"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string `mapstructure:"secret" validate:"required,min=32"` // 安全起见，至少32字符
	Issuer string `mapstructure:"issuer" validate:"required"`
	Expire int64  `mapstructure:"expire" validate:"required,min=60"` // 单位秒
}

// StorageConfig 存储配置
type StorageConfig struct {
	Backend string              `mapstructure:"backend" validate:"required,oneof=local s3 webdav"`
	Access  StorageAccessConfig `mapstructure:"access" validate:"required"`

	Local  *LocalStorageConfig  `mapstructure:"local" validate:"required_if=Backend local"`
	S3     *S3StorageConfig     `mapstructure:"s3" validate:"required_if=Backend s3"`
	WebDAV *WebDAVStorageConfig `mapstructure:"webdav" validate:"required_if=Backend webdav"`
}

// StorageAccessConfig 访问策略
type StorageAccessConfig struct {
	GinProxy   GinProxyConfig    `mapstructure:"gin_proxy"`
	ImageProxy *ImageProxyConfig `mapstructure:"image_proxy"`
}

// GinProxyConfig 本地 HTTP 代理访问
type GinProxyConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// ImageProxyConfig 图片代理（如 imgproxy / thumbor）
type ImageProxyConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	BaseURL string `mapstructure:"base_url" validate:"required_if=Enabled true,url"`
}

// LocalStorageConfig 本地存储
type LocalStorageConfig struct {
	Root string `mapstructure:"root" validate:"required"`
}

// S3StorageConfig S3 / S3-compatible 存储
type S3StorageConfig struct {
	UseSSL   bool   `mapstructure:"use_ssl"`
	Endpoint string `mapstructure:"endpoint" validate:"required"`
	Region   string `mapstructure:"region"`
	Bucket   string `mapstructure:"bucket" validate:"required"`

	Credentials struct {
		AccessKeyID     string `mapstructure:"access_key_id" validate:"required"`
		SecretAccessKey string `mapstructure:"secret_access_key" validate:"required"`
	} `mapstructure:"credentials" validate:"required"`

	PublicBaseURL  string `mapstructure:"public_base_url" validate:"required,url"`
	PresignEnabled bool   `mapstructure:"presign_enabled"`
	PresignExpire  int64  `mapstructure:"presign_expire" validate:"omitempty,min=1"`
}

// WebDAVStorageConfig WebDAV 存储
type WebDAVStorageConfig struct {
	Endpoint      string `mapstructure:"endpoint" validate:"required,url"`
	BasePath      string `mapstructure:"base_path" validate:"required"`
	PublicBaseURL string `mapstructure:"public_base_url" validate:"required,url"`

	Auth struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"auth"`
}
