package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	AdminToken       string           `yaml:"adminToken"`
	Kubernetes       KubernetesConfig `yaml:"kubernetes"`
	Redis            RedisConfig      `yaml:"redis"`
	LogLevel         string           `yaml:"logLevel"`
	InstructionsPath string           `yaml:"instructionsPath"`
}

type RedisConfig struct {
	Address   string `yaml:"address"`
	IsCluster bool   `yaml:"isCluster"`
}

type KubernetesConfig struct {
	DevMode bool `yaml:"devMode"`
	// this must be set if dev mode is set to true
	Kubeconfig string `yaml:"kubeconfig"`
	ApiURI     string `yaml:"apiUri"`
	ApiCA      string `yaml:"apiCa"`
}

var (
	doOnce   sync.Once
	v        *viper.Viper
	cfg      *Config
	updates  chan Config
	kupdates chan KubernetesConfig
)

func Get() (*Config, <-chan Config) {
	doOnce.Do(func() {
		cfg = &Config{}
		updates = make(chan Config)
		kupdates = make(chan KubernetesConfig)
		defaultAdminToken := uuid.Must(uuid.NewV4()).String()
		v = viper.New()
		v.SetEnvPrefix("ksk")
		v.SetConfigType("yaml")
		v.SetConfigName("config")
		v.AddConfigPath(".")

		v.SetDefault("redis.isCluster", false)
		v.SetDefault("redis.address", "redis:6379")
		v.SetDefault("logLevel", "info")
		v.SetDefault("instructionsPath", "/ksk/instructions.yaml")
		v.SetDefault("adminToken", defaultAdminToken)
		v.SetDefault("kubernetes.apiUri", "")
		v.SetDefault("kubernetes.apiCa", "")

		log.Info("Default admin token: ", defaultAdminToken)

		v.AutomaticEnv()
		if err := v.ReadInConfig(); err != nil {
			log.WithError(err).Error("failed to read Config")
		} else {
			v.WatchConfig()
			v.OnConfigChange(notify)
		}
		err := v.Unmarshal(cfg)
		if err != nil {
			panic(err)
		}
	})
	return cfg, updates
}

func GetKubernetesConfig() (KubernetesConfig, <-chan KubernetesConfig) {
	cfg, _ := Get()
	return cfg.Kubernetes, kupdates
}

func GetRedisConfig() RedisConfig {
	cfg, _ := Get()
	return cfg.Redis
}

func notify(e fsnotify.Event) {
	log.WithField("filename", e.Name).Info("Config update")
	err := v.Unmarshal(cfg)
	if err != nil {
		log.WithError(err).Error("failed to unmarshal Config in notify")
		return
	}
	cfg.AdminToken = v.GetString("adminToken")
	select {
	// do not care if nobody is reading
	case updates <- *cfg:
	default:
		break
	}
	select {
	// do not care if nobody is reading
	case kupdates <- cfg.Kubernetes:
	default:
		break
	}
}
