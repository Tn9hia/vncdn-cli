package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// LoadConfig loads the CDN configuration from a file.
type Profile struct {
    Name            string `mapstructure:"name"`
    AccessKey       string `mapstructure:"accessKey"`
    AccessKeySecret string `mapstructure:"accessKeySecret"`
}

type Profiles struct {
    DefaultProfile string    `mapstructure:"default_profile"`
    Profiles       []Profile `mapstructure:"profiles"`
}

// InitializeConfig sets up Viper and checks/creates the config file
func initializeConfig() error {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	configPath := filepath.Join(homeDir, ".config", "cdnctl")
	configFile := filepath.Join(configPath, "config.yaml")

	// Set the configuration file and type
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	// Check if the config file exists, if not create it with 2 fields (default_profile and profiles)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		// Create a new config file with default values
		viper.Set("default_profile", "")
		viper.Set("profiles", []Profile{})
		if err := viper.WriteConfigAs(configFile); err != nil {
			return fmt.Errorf("failed to write default config file: %w", err)
		}
	}

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	return nil
}

// GetDefaultProfile returns the credentials for the default profile
func GetDefaultProfile() (string, string, error) {
	var setting Profiles

	// Initialize the configuration
	if err := initializeConfig(); err != nil {
		return "", "", fmt.Errorf("error initializing config: %w", err)
	}

	// Load profiles from config
	if err := viper.Unmarshal(&setting); err != nil {
		return "", "", fmt.Errorf("error unmarshalling profiles: %w", err)
	}

	if len(setting.Profiles) == 0 {
		return "", "", fmt.Errorf("no profiles found. Please add a profile using 'cdnctl config add'")
	}

	// If no default profile is set, use the first one
	defaultProfileName := setting.DefaultProfile
	if defaultProfileName == "" && len(setting.Profiles) > 0 {
		defaultProfileName = setting.Profiles[0].Name
	}

	// Find the default profile
	for _, profile := range setting.Profiles {
		if profile.Name == defaultProfileName {
			return profile.AccessKey, profile.AccessKeySecret, nil
		}
	}

	return "", "", fmt.Errorf("default profile '%s' not found", defaultProfileName)
}

// Retrieve and display profiles
func DisplayProfiles(profileName string, isRaw bool) (string, string, error) {
	var setting Profiles

	// Initialize the configuration
	if err := initializeConfig(); err != nil {
		return "", "", fmt.Errorf("error initializing config: %w", err)
	}

	// Load profiles from config
	if err := viper.Unmarshal(&setting); err != nil {
		fmt.Printf("error unmarshalling profiles: %v", err)
		return "", "", err
	}

	if len(setting.Profiles) == 0 {
		if !isRaw {
			fmt.Println("No profiles found. Please add a profile using 'cdnctl config add'.")
		}
		return "", "", fmt.Errorf("no profiles found")
	}

	// Case: profileName được chỉ định
	if profileName != "" {
		for _, profile := range setting.Profiles {
			if profile.Name == profileName {
				if isRaw {
					// Trả về access key và secret key
					return profile.AccessKey, profile.AccessKeySecret, nil
				} else {
					// In ra như bình thường
					fmt.Printf("Profile Name: %s\nAccess Key: %s\nAccess Key Secret: %s\n",
						profile.Name, profile.AccessKey, profile.AccessKeySecret)
					return "", "", nil
				}
			}
		}
		if !isRaw {
			fmt.Printf("Profile '%s' not found.\n", profileName)
		}
		return "", "", fmt.Errorf("profile '%s' not found", profileName)
	}

	// Nếu không có profileName và không phải raw thì in tất cả
	if !isRaw {
		fmt.Printf("Default Profile: %s\n\n", setting.DefaultProfile)
		fmt.Println("Available CDN Profiles:")
		for _, profile := range setting.Profiles {
			marker := ""
			if profile.Name == setting.DefaultProfile {
				marker = " (default)"
			}
			fmt.Printf("Name: %s%s\n   AccessKey: %s\n   SecretKey: %s\n\n",
				profile.Name, marker, profile.AccessKey, profile.AccessKeySecret)
		}
	}

	return "", "", nil
}

// Add profile
func AddProfile(name, accessKey, accessKeySecret, defaultProfile string) error {
	var setting Profiles
	
	// Initialize the configuration
	if err := initializeConfig(); err != nil {
		return fmt.Errorf("error initializing config: %w", err)
	}
	// Read the existing profiles from the config file
	if err := viper.Unmarshal(&setting); err != nil {
		return fmt.Errorf("error unmarshalling profile: %w", err)
	}
	// Check if profile already exists
	for _, profile := range setting.Profiles {
		if profile.Name == name {
			return fmt.Errorf("profile with name '%s' already exists", name)
		}
	}

	// Append the new profile to the profiles slice
	setting.Profiles = append(setting.Profiles, Profile{
		Name:            name,
		AccessKey:       accessKey,
		AccessKeySecret: accessKeySecret,
	})

	// Set the updated profiles back to Viper
	viper.Set("profiles", setting.Profiles)

	// Check if default profile is empty and set it if necessary
	if setting.DefaultProfile == "" || defaultProfile == "Yes" || defaultProfile == "yes" {
		setting.DefaultProfile = name
		viper.Set("default_profile", setting.DefaultProfile)
	}

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	
	fmt.Printf("Profile '%s' added successfully.\n", name)
	if setting.DefaultProfile == name {
		fmt.Printf("Set as default profile.\n")
	}
	return nil
}

// Remove profile
func RemoveProfile(name string) error {
	var profile Profiles
	// Initialize the configuration
	if err := initializeConfig(); err != nil {
		return fmt.Errorf("error initializing config: %w", err)
	}
	// Read the existing profiles from the config file
	if err := viper.Unmarshal(&profile); err != nil {
		return fmt.Errorf("error unmarshalling profile: %w", err)
	}

	// Find and remove the profile
	for i, p := range profile.Profiles {
		if p.Name == name {
			profile.Profiles = append(profile.Profiles[:i], profile.Profiles[i+1:]...)
			
			// If the removed profile was the default, clear default or set to first available
			if profile.DefaultProfile == name {
				if len(profile.Profiles) > 0 {
					profile.DefaultProfile = profile.Profiles[0].Name
					fmt.Printf("Default profile changed to '%s'.\n", profile.DefaultProfile)
				} else {
					profile.DefaultProfile = ""
				}
				viper.Set("default_profile", profile.DefaultProfile)
			}
			
			viper.Set("profiles", profile.Profiles)
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("error writing config file: %w", err)
			}
			fmt.Printf("Profile '%s' removed successfully.\n", name)
			return nil
		}
	}

	return fmt.Errorf("profile with name '%s' not found", name)
}