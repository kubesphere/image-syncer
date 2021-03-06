package client

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

// Config information of sync client
type Config struct {
	// the authentication information of each registry
	AuthList map[string]Auth `json:"auth" yaml:"auth"`

	// a <source_repo>:<dest_repo> map
	ImageList map[string]string `json:"images" yaml:"images"`

	// If the destinate registry and namespace is not provided,
	// the source image will be synchronized to defaultDestRegistry
	// and defaultDestNamespace with origin repo name and tag.
	defaultDestRegistry  string
	defaultDestNamespace string
}

// Auth describes the authentication information of a registry
type Auth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Insecure bool   `json:"insecure" yaml:"insecure"`
}

// NewSyncConfig creates a Config struct
func NewSyncConfig(authFilePath, imageFilePath, defaultDestRegistry, defaultDestNamespace string) (*Config, error) {
	config := Config{
		defaultDestNamespace: defaultDestNamespace,
		defaultDestRegistry:  defaultDestRegistry,
	}

	if err := openAndDecode(authFilePath, &config.AuthList); err != nil {
		return nil, fmt.Errorf("decode auth file %v error: %v", authFilePath, err)
	}

	if err := openAndDecode(imageFilePath, &config.ImageList); err != nil {
		return nil, fmt.Errorf("decode image file %v error: %v", imageFilePath, err)
	}

	return &config, nil
}

// Open json file and decode into target interface
func openAndDecode(filePath string, target interface{}) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %v not exist: %v", filePath, err)
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("open file %v error: %v", filePath, err)
	}

	if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(target); err != nil {
			return fmt.Errorf("unmarshal config error: %v", err)
		}
	} else if strings.HasSuffix(filePath, ".json") {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(target); err != nil {
			return fmt.Errorf("unmarshal config error: %v", err)
		}
	} else {
		return fmt.Errorf("only one of yaml/yml/json format is supported")
	}

	return nil
}

// GetAuth gets the authentication information in Config
func (c *Config) GetAuth(registry string, namespace string) (Auth, bool) {
	// key of each AuthList item can be "registry/namespace" or "registry" only
	registryAndNamespace := registry + "/" + namespace

	if moreSpecificAuth, exist := c.AuthList[registryAndNamespace]; exist {
		return moreSpecificAuth, exist
	}

	auth, exist := c.AuthList[registry]
	return auth, exist
}

// GetImageList gets the ImageList map in Config
func (c *Config) GetImageList() map[string]string {
	return c.ImageList
}
