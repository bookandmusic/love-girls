package config

// AppConfig 根配置
type AppConfig struct {
	App        AppConfigApp     `mapstructure:"app" validate:"required"`
	Server     ServerConfig     `mapstructure:"server" validate:"required"`
	Log        LogConfig        `mapstructure:"log" validate:"required"`
	DataSource DataSourceConfig `mapstructure:"datasource" validate:"required"`
	JWT        JWTConfig        `mapstructure:"jwt" validate:"required"`
	Storage    StorageConfig    `mapstructure:"storage" validate:"required"`
	ImageProxy ImageProxyConfig `mapstructure:"image_proxy"`
}

// AppConfigApp 应用基本信息
type AppConfigApp struct {
	Name string `mapstructure:"name" validate:"required"`
	Env  string `mapstructure:"env" validate:"required,oneof=dev test prod"`
}

// ServerConfig HTTP 服务配置
type ServerConfig struct {
	Addr        string `mapstructure:"addr" validate:"required"`
	Mode        string `mapstructure:"mode" validate:"required,oneof=debug release test"`
	InternalURL string `mapstructure:"internal_url"` // 内网地址，ImageProxy 访问 Gin 用（可选）
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
	Local   *LocalStorageConfig `mapstructure:"local" validate:"required_if=Backend local"`
	S3      *S3StorageConfig    `mapstructure:"s3" validate:"required_if=Backend s3"`
	WebDAV  *WebDAVStorageConfig `mapstructure:"webdav" validate:"required_if=Backend webdav"`
}

// ImageProxyConfig 图片代理（如 imgproxy / thumbor）
type ImageProxyConfig struct {
	InternalURL string `mapstructure:"internal_url"` // 内网地址，Gin 转发时使用
	PublicURL   string `mapstructure:"public_url"`   // 公开地址，前端直接访问（可选）
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

	PublicURL     string `mapstructure:"public_url"`     // 公开访问地址（可选）
	PresignEnable bool   `mapstructure:"presign_enable"` // 是否启用预签名
	PresignExpire int64  `mapstructure:"presign_expire" validate:"omitempty,min=1"`
}

// WebDAVStorageConfig WebDAV 存储
type WebDAVStorageConfig struct {
	Endpoint  string `mapstructure:"endpoint" validate:"required,url"`
	BasePath  string `mapstructure:"base_path" validate:"required"`
	PublicURL string `mapstructure:"public_url"` // 公开访问地址（可选）

	Auth struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"auth"`
}
