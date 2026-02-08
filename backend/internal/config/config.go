package config

import (
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
}

// LoadConfig 加载应用配置（线程安全）
func LoadConfig() (*AppConfig, error) {
	var err error
	once.Do(func() {
		_ = godotenv.Load() // 优先加载 .env

		v := viper.New()
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./data/configs")
		v.AddConfigPath(".")

		// Step 1: 默认值注入
		setDefaults(v)

		// Step 3: 加载与环境变量覆盖（AutomaticEnv 优先于 ReadInConfig）
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		readErr := v.ReadInConfig()

		var c AppConfig

		if readErr != nil {
			// Step 2: 文件检测与创建
			configDir := "./data/configs"
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

	v.SetDefault("server.addr", ":8182")
	v.SetDefault("server.schema", "http")
	v.SetDefault("server.host_name", "localhost:8182")
	v.SetDefault("server.mode", "debug")

	v.SetDefault("log.level", "debug")

	v.SetDefault("jwt.secret", "love-girl-123456789012345678901234")
	v.SetDefault("jwt.issuer", "love-girl")
	v.SetDefault("jwt.expire", 900)

	v.SetDefault("datasource.database.driver", "sqlite")
	v.SetDefault("datasource.database.dsn", "./data/love-girl.db")

	v.SetDefault("storage.backend", "local")
	v.SetDefault("storage.access.gin_proxy.enabled", true)
	v.SetDefault("storage.access.image_proxy.enabled", false)
	v.SetDefault("storage.access.image_proxy.base_url", "http://localhost:8080")
	v.SetDefault("storage.local.root", "./data/uploads")

	// S3 存储：仅注册 key（不设默认值，避免创建空结构体触发内部 required 校验）
	_ = v.BindEnv("storage.s3.use_ssl")
	_ = v.BindEnv("storage.s3.endpoint")
	_ = v.BindEnv("storage.s3.region")
	_ = v.BindEnv("storage.s3.bucket")
	_ = v.BindEnv("storage.s3.credentials.access_key_id")
	_ = v.BindEnv("storage.s3.credentials.secret_access_key")
	_ = v.BindEnv("storage.s3.public_base_url")
	_ = v.BindEnv("storage.s3.presign_enabled")
	_ = v.BindEnv("storage.s3.presign_expire")

	// WebDAV 存储：仅注册 key
	_ = v.BindEnv("storage.webdav.endpoint")
	_ = v.BindEnv("storage.webdav.base_path")
	_ = v.BindEnv("storage.webdav.public_base_url")
	_ = v.BindEnv("storage.webdav.auth.username")
	_ = v.BindEnv("storage.webdav.auth.password")
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
	configDir := "./data/configs"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config directory failed: %w", err)
	}

	configPath := configDir + "/config.yaml"

	// 读取现有配置
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)

	setDefaults(v)
	_ = v.ReadInConfig()

	// 设置指定的键值
	v.Set(key, value)

	// 写入配置文件
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
