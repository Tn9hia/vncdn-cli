package cmd

import (
	"cdnctl/utils"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage token for cdn profiles",
}

var addConfig = &cobra.Command{
	Use:   "add",
	Short: "Add a new CDN profile configuration",
	Long: `Use this command to add a new CDN profile configuration. You will need to provide the necessary details such as profile name, token, and other relevant information.`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt := promptui.Prompt{
			Label: "Enter Profile Name",
		}
		profileName, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		prompt = promptui.Prompt{
			Label: "Enter Access Key",
		}
		accessKey, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		prompt = promptui.Prompt{
			Label: "Enter Access Key Secret",
			Mask: '*', // Use Mask instead of HideEntered
		}
		accessKeySecret, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		// If set as default profile, it will be set as default
		promptSelect := promptui.Select{
			Label: "Set as Default Profile?",
			Items: []string{"Yes", "No"},
		}
		_, defaultProfile, err := promptSelect.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = utils.AddProfile(profileName, accessKey, accessKeySecret, defaultProfile)
		if err != nil {
			fmt.Printf("Error adding profile: %v\n", err)
		}
	},
}

var removeConfig = &cobra.Command{
	Use:   "remove",
	Short: "Remove an existing CDN profile configuration",
	Long: `Use this command to remove an existing CDN profile configuration. You will need to specify the profile name or ID of the configuration you want to remove.`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt := promptui.Prompt{
			Label: "Enter Profile Name to Remove",
		}
		profileName, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		err = utils.RemoveProfile(profileName)
		if err != nil {
			fmt.Printf("Error removing profile: %v\n", err)
		}
	},
}

var viewConfig = &cobra.Command{
	Use:   "show",
	Short: "View existing CDN profile configurations",
	Long: `Use this command to view all existing CDN profile configurations. It will display the details of each configuration, such as profile name, token, and other relevant information.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// If a profile name is provided, display that specific profile
			profileName := args[0]
			utils.DisplayProfiles(profileName, false)
			return
		}
		utils.DisplayProfiles("", false) // Display all profiles if no specific profile is provided
	},
}

func init() {
	// Add the config command to the root command
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(addConfig)
	configCmd.AddCommand(removeConfig)
	configCmd.AddCommand(viewConfig)

	// REMOVED: The duplicate rootCmd initialization and execution that was causing issues
}