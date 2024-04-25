package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/go-viper/mapstructure/v2"
	"github.com/google/uuid"

	"github.com/alexZaicev/message-broker/internal/domain/entities"
)

const (
	configPath       = "MB_CONF"
	defaultConfigDir = ".mb"
)

func LoadConfiguration(mbCtx *entities.MessageBrokerContext) error {
	resources, err := listResources()
	if err != nil {
		return err
	}

	for _, resource := range resources {
		switch resource.Kind {
		case ResourceKindQueue:
			var queueSpec QueueSpec
			if decodeErr := mapstructure.Decode(resource.Spec, &queueSpec); decodeErr != nil {
				return fmt.Errorf("failed to decode queue spec: %w", decodeErr)
			}

			mbCtx.AddExchange(entities.NewQueueExchange(uuid.New(), queueSpec.Name))
		case ResourceKindTopic:
			var topicSpec TopicSpec
			if decodeErr := mapstructure.Decode(resource.Spec, &topicSpec); decodeErr != nil {
				return fmt.Errorf("failed to decode topic spec: %w", decodeErr)
			}

			mbCtx.AddExchange(entities.NewTopicExchange(uuid.New(), topicSpec.Name, topicSpec.Pattern))
		}
	}

	return nil
}

func listResources() ([]Resource, error) {
	path, err := getConfigAbsPath()
	if err != nil {
		return nil, fmt.Errorf("failed to obtain path to configuration files: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to gather information of the configuration directory: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("provided configuration path is not a directory: %s", path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list configuration directory: %w", err)
	}

	var resources []Resource

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		bytes, readErr := os.ReadFile(filepath.Join(path, entry.Name()))
		if readErr != nil {
			return nil, fmt.Errorf("failed to read configuration file: %w", readErr)
		}

		var resource Resource
		if yamlErr := yaml.Unmarshal(bytes, &resource); yamlErr != nil {
			return nil, fmt.Errorf("failed to unmarshal configuration file: %w", yamlErr)
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func getConfigAbsPath() (string, error) {
	val := os.Getenv(configPath)
	if val != "" {
		absPath, err := filepath.Abs(val)
		if err != nil {
			return "", err
		}

		return absPath, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, defaultConfigDir), nil
}
