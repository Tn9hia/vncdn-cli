package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"cdnctl/utils"
	"encoding/json"
)

// wsaCmd represents the wsa command
var wsaCmd = &cobra.Command{
	Use:   "wsa",
	Short: "Manage Web Acceleration Services",
}

var wsaGetCmd = &cobra.Command{
	Use:   "get [domain]",
	Short: "Get Web Acceleration Services for a domain",
	Args:  cobra.ExactArgs(1), // Require exactly one argument (domain)
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0] // Get domain from command line arguments
		
		// Create request body with the domain from command line
		body := map[string]interface{}{
			"domain": domain,
		}
		
		// Convert body to JSON
		bodyJson, err := json.Marshal(body)
		if err != nil {
			fmt.Printf("Error marshaling request body: %v\n", err)
			return
		}
		
		// Call the API with proper error handling
		data, err := utils.CallApi("POST", utils.BaseURL1+"/v1.1/service_id", "/v1.1/service_id", bodyJson, "")
		if err != nil {
			fmt.Printf("API call failed: %v\n", err)
			return
		}
		
		// Pretty print the JSON response
		var prettyJSON map[string]interface{}
		if err := json.Unmarshal(data, &prettyJSON); err != nil {
			fmt.Printf("Response: %s\n", string(data))
		} else {
			prettyData, _ := json.MarshalIndent(prettyJSON, "", "  ")
			fmt.Printf("Response:\n%s\n", string(prettyData))
		}
	},
}

func init() {
	rootCmd.AddCommand(wsaCmd)
	wsaCmd.AddCommand(wsaGetCmd)
}