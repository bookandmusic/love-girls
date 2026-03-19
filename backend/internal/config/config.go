package config

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	cfg       *AppConfig
	once      sync.Once
	viperInst *viper.Viper // 保存 viper 实例用于配置监听
)

// validator 实例
var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("mapstructure")
		if name != "" && name != "-" {
			return name
		}
		return fld.Name
	})

	// 注册结构体级别验证器
	validate.RegisterStructValidation(validateStorageConfig, StorageConfig{})
}

// LoadConfig 加载应用配置（线程安全）
func LoadConfig() (*AppConfig, error) {
	var err error
	once.Do(func() {
		_ = godotenv.Load() // 优先加载 .env

		// 先确定 data_dir（从环境变量或使用默认值）
		dataDir := os.Getenv("DATA_DIR")
		if dataDir == "" {
			dataDir = "./data"
		}
		configDir := dataDir + "/configs"

		v := viper.New()
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(configDir)

		// Step 1: 默认值注入
		setDefaults(v)

		// Step 3: 加载与环境变量覆盖（AutomaticEnv 优先于 ReadInConfig）
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		readErr := v.ReadInConfig()

		var c AppConfig

		if readErr != nil {
			// Step 2: 文件检测与创建
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return
			}

			configPath := configDir + "/config.yaml"
			// 使用 SafeWriteConfigAs 确保不覆盖现有文件
			if writeErr := v.SafeWriteConfigAs(configPath); writeErr != nil {
				err = fmt.Errorf("write default config failed: %w", writeErr)
				return
			}

			// 重新读取配置文件
			if readErr := v.ReadInConfig(); readErr != nil {
				err = fmt.Errorf("read config after write failed: %w", readErr)
				return
			}
		}

		if e := v.Unmarshal(&c); e != nil {
			err = fmt.Errorf("unmarshal config failed: %w", e)
			return
		}

		// 根据 data_dir 计算默认路径
		applyDataDirDefaults(&c, v)

		// 统一创建数据目录
		if e := initDataDir(&c); e != nil {
			err = fmt.Errorf("init data directory failed: %w", e)
			return
		}

		// 自动生成 JWT Secret（如果未配置）并持久化
		applyJWTDefaults(&c, v)

		if e := validate.Struct(&c); e != nil {
			err = e
			return
		}

		ApplyEnvPolicy(&c)

		cfg = &c
		viperInst = v // 保存 viper 实例
	})

	return cfg, err
}

// GetViperInstance 获取 viper 实例用于配置监听
func GetViperInstance() *viper.Viper {
	return viperInst
}

// ResetConfig 重置配置状态，使 LoadConfig 可被再次调用
func ResetConfig() {
	once = sync.Once{}
	cfg = nil
	viperInst = nil
}

// setDefaults 集中管理默认值
func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "love-girl")
	v.SetDefault("app.env", "prod")
	v.SetDefault("data_dir", "./data") // 统一数据目录

	v.SetDefault("server.addr", ":8182")
	v.SetDefault("server.mode", "debug")

	// CORS 配置：默认允许所有来源，常用方法和请求头
	v.SetDefault("server.cors.allow_all_origins", true)
	v.SetDefault("server.cors.allow_origins", []string{})
	v.SetDefault("server.cors.allow_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"})
	v.SetDefault("server.cors.allow_headers", []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"})

	// CORS 环境变量绑定
	_ = v.BindEnv("server.cors.allow_all_origins", "CORS_ALLOW_ALL_ORIGINS")
	_ = v.BindEnv("server.cors.allow_origins", "CORS_ALLOW_ORIGINS")

	v.SetDefault("log.level", "debug")
	v.SetDefault("log.format", "text")
	v.SetDefault("log.output", "stdout")

	// JWT 配置：issuer 和 expire 有默认值，secret 由程序自动生成
	v.SetDefault("jwt.issuer", "love-girl")
	v.SetDefault("jwt.expire", 86400) // 24小时
	_ = v.BindEnv("jwt.secret", "JWT_SECRET")
	_ = v.BindEnv("jwt.issuer", "JWT_ISSUER")
	_ = v.BindEnv("jwt.expire", "JWT_EXPIRE")

	// 数据库：默认使用 data_dir 下的 SQLite
	v.SetDefault("datasource.database.driver", "sqlite")
	// DSN 不设默认值，由 applyDataDirDefaults 根据 data_dir 计算

	// 存储：默认使用 data_dir 下的 local 存储
	v.SetDefault("storage.backend", "local")
	// Local 存储路径由 data_dir 自动计算，不支持配置

	v.SetDefault("image_proxy.internal_url", "")
	v.SetDefault("image_proxy.public_url", "")

	// 环境变量绑定
	_ = v.BindEnv("data_dir", "DATA_DIR")
	_ = v.BindEnv("datasource.database.driver", "DATABASE_DRIVER")
	_ = v.BindEnv("datasource.database.dsn", "DATABASE_DSN")
	_ = v.BindEnv("storage.backend", "STORAGE_BACKEND")

	// S3 存储：仅注册 key（不设默认值，避免创建空结构体触发内部 required 校验）
	_ = v.BindEnv("storage.s3.use_ssl")
	_ = v.BindEnv("storage.s3.endpoint")
	_ = v.BindEnv("storage.s3.region")
	_ = v.BindEnv("storage.s3.bucket")
	_ = v.BindEnv("storage.s3.credentials.access_key_id")
	_ = v.BindEnv("storage.s3.credentials.secret_access_key")
	_ = v.BindEnv("storage.s3.public_url")
	_ = v.BindEnv("storage.s3.presign_enable")
	_ = v.BindEnv("storage.s3.presign_expire")

	// WebDAV 存储：仅注册 key
	_ = v.BindEnv("storage.webdav.endpoint")
	_ = v.BindEnv("storage.webdav.base_path")
	_ = v.BindEnv("storage.webdav.public_url")
	_ = v.BindEnv("storage.webdav.auth.username")
	_ = v.BindEnv("storage.webdav.auth.password")

	// ImageProxy：仅注册 key
	_ = v.BindEnv("image_proxy.internal_url")
	_ = v.BindEnv("image_proxy.public_url")
}

// validateStorageConfig 结构体级别验证：根据 Backend 类型验证对应配置
func validateStorageConfig(sl validator.StructLevel) {
	config := sl.Current().Interface().(StorageConfig)

	switch config.Backend {
	case "s3":
		if config.S3 == nil {
			sl.ReportError(sl.Current(), "S3", "S3", "required", "")
			return
		}
		if config.S3.Endpoint == "" {
			sl.ReportError(reflect.ValueOf(config.S3), "Endpoint", "endpoint", "required", "")
		}
		if config.S3.Bucket == "" {
			sl.ReportError(reflect.ValueOf(config.S3), "Bucket", "bucket", "required", "")
		}
		if config.S3.Credentials.AccessKeyID == "" {
			sl.ReportError(reflect.ValueOf(config.S3), "Credentials.AccessKeyID", "credentials.access_key_id", "required", "")
		}
		if config.S3.Credentials.SecretAccessKey == "" {
			sl.ReportError(reflect.ValueOf(config.S3), "Credentials.SecretAccessKey", "credentials.secret_access_key", "required", "")
		}

	case "webdav":
		if config.WebDAV == nil {
			sl.ReportError(sl.Current(), "WebDAV", "WebDAV", "required", "")
			return
		}
		if config.WebDAV.Endpoint == "" {
			sl.ReportError(reflect.ValueOf(config.WebDAV), "Endpoint", "endpoint", "required", "")
		}
		if config.WebDAV.BasePath == "" {
			sl.ReportError(reflect.ValueOf(config.WebDAV), "BasePath", "base_path", "required", "")
		}
	}
}

// applyDataDirDefaults 根据 data_dir 计算默认路径
func applyDataDirDefaults(cfg *AppConfig, v *viper.Viper) {
	paths := cfg.GetDataPaths()

	// SQLite 数据库：DSN 由 data_dir 自动计算，忽略用户配置
	if cfg.DataSource.Database.Driver == "sqlite" {
		cfg.DataSource.Database.DSN = paths.DBFile
		v.Set("datasource.database.dsn", paths.DBFile)
	}

	// MySQL/PostgreSQL：DSN 必须由用户配置
}

// initDataDir 统一创建数据目录
func initDataDir(cfg *AppConfig) error {
	paths := cfg.GetDataPaths()

	// 创建数据根目录
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}

	// 创建配置目录
	if err := os.MkdirAll(paths.ConfigDir, 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	// 创建上传目录（仅 local 存储需要）
	if cfg.Storage.Backend == "local" {
		if err := os.MkdirAll(paths.UploadDir, 0755); err != nil {
			return fmt.Errorf("create upload directory: %w", err)
		}
	}

	// 创建数据库目录（SQLite 需要）
	if cfg.DataSource.Database.Driver == "sqlite" {
		dbDir := filepath.Dir(cfg.DataSource.Database.DSN)
		if dbDir != "" && dbDir != "." {
			if err := os.MkdirAll(dbDir, 0755); err != nil {
				return fmt.Errorf("create database directory: %w", err)
			}
		}
	}

	return nil
}

// applyJWTDefaults 自动生成 JWT Secret（如果未配置）并持久化到配置文件
func applyJWTDefaults(cfg *AppConfig, v *viper.Viper) {
	if cfg.JWT.Secret == "" {
		// 生成 64 字符随机密钥
		bytes := make([]byte, 32)
		if _, err := cryptorand.Read(bytes); err != nil {
			// 如果随机生成失败，使用备用方案
			cfg.JWT.Secret = "fallback-secret-" + fmt.Sprintf("%d", os.Getpid())
		} else {
			cfg.JWT.Secret = hex.EncodeToString(bytes)
		}
		cfg.JWT._autoGenerated = true

		// 持久化到配置文件
		v.Set("jwt.secret", cfg.JWT.Secret)
		paths := cfg.GetDataPaths()
		_ = v.WriteConfigAs(paths.ConfigDir + "/config.yaml")
	}
}

// ApplyEnvPolicy 根据环境设置运行策略
func ApplyEnvPolicy(cfg *AppConfig) {
	switch cfg.App.Env {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
		if cfg.Log.Level == "debug" {
			cfg.Log.Level = "info"
		}
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

// UpdateConfigValue 更新配置文件中的指定键值
func UpdateConfigValue(key string, value any) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	paths := cfg.GetDataPaths()

	// 读取现有配置
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(paths.ConfigDir)

	setDefaults(v)
	_ = v.ReadInConfig()

	// 设置指定的键值
	v.Set(key, value)

	// 写入配置文件
	configPath := paths.ConfigDir + "/config.yaml"
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("write config file failed: %w", err)
	}

	return nil
}

// WatchConfig 监听配置文件变化，校验通过后通过 channel 通知调用方重启
func WatchConfig(v *viper.Viper, configPath string, restartCh chan<- struct{}) {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if e.Op != fsnotify.Write || filepath.Base(e.Name) != filepath.Base(configPath) {
			return
		}

		// 原子性校验逻辑：创建新的 Viper 实例尝试解析配置
		vTest := viper.New()
		vTest.SetConfigName("config")
		vTest.SetConfigType("yaml")
		vTest.AddConfigPath(filepath.Dir(configPath))

		setDefaults(vTest)
		vTest.AutomaticEnv()
		vTest.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if readErr := vTest.ReadInConfig(); readErr != nil {
			fmt.Printf("配置文件读取失败: %v\n", readErr)
			return
		}

		var testCfg AppConfig
		if unmarshalErr := vTest.Unmarshal(&testCfg); unmarshalErr != nil {
			fmt.Printf("配置文件解析失败（语法错误）: %v\n", unmarshalErr)
			return
		}

		if validateErr := validate.Struct(&testCfg); validateErr != nil {
			fmt.Printf("配置验证失败（业务规则错误）: %v\n", validateErr)
			return
		}

		fmt.Println("配置验证通过，正在触发进程内重启...")

		select {
		case restartCh <- struct{}{}:
		default:
			// channel 已满，说明已有重启信号待处理
		}
	})
}
