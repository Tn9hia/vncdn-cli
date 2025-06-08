package cmd

import (
	"fmt"
	// "strconv"
	"github.com/spf13/cobra"
	"cdnctl/utils"
	"encoding/json"
)

// wsaCmd represents the wsa command
var wsaCmd = &cobra.Command{
	Use:   "wsa",
	Short: "Manage Web Accessleration Services",
}

var wsaGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Web Accessleration Services",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("wsa called")
		// Temporary command 
		body := map[string]interface{}{
			"domain": "h2.riok.site",
		}
		var data []byte
		var err error
		bodyJson, _ := json.Marshal(body) // Add JSON encoding
		data, err = utils.CallApi("POST", utils.BaseURL1+"/v1.1/service_id", string(bodyJson), nil, "") // Update function call
		if err != nil {
			panic(err)
		}
		fmt.Println("Response received from WSA API:", string(data))
		fmt.Println("Response:", string(data))
		
	},
}

func init() {
	rootCmd.AddCommand(wsaCmd)
	wsaCmd.AddCommand(wsaGetCmd)

}
