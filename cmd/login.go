/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"cookiefieldcli/cmd/login"
	// "cookiefieldcli/cmd/login/interface"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with the cli.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Viper calculates the path of the config file starting from where go.mod is located
		viper.AddConfigPath("./configs")
		viper.SetConfigName("config") // Register config file name (no extension)
		viper.SetConfigType("json")   // Look for specific type
		viper.ReadInConfig()
		domain := viper.Get("domain")
		clientId := viper.Get("client_id")
		fmt.Println("login called")
		// fmt.Println(domain)
		// fmt.Println(clientId)
		loginJob := login.NewLoginJob(clientId.(string), domain.(string))
		//First we need to retreive the device code.
		deviceCodeErr := loginJob.GetDeviceCode()
		if deviceCodeErr != nil {
			log.Panic("Requesting Device Code failed.", deviceCodeErr)
		}
		// fmt.Println("Response from GetDeviceCode: ", loginJob.DeviceCodeData)
		//Second, we need to get a request token.
		ReqTokenErr := loginJob.GetRequestToken()
		if ReqTokenErr != nil {
			log.Panic("Requesting Token failed.", ReqTokenErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
