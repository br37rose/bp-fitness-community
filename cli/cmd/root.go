package cmd

import (
	"fmt"
	"os"

	// homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var (
// databaseHost                      string
// databasePort                      string
// databaseUser                      string
// databasePassword                  string
// databaseName                      string
// databaseSSLMode                   string
// applicationIP                     string
// applicationPort                   string
// applicationSigningKey             string
// invoicebuilderPdfTemplateFilePath string
// invoicebuilderData                string
// hasAutoMigrations                 bool
)

// Initialize function will be called when every command gets called.
func init() {
	// // Get our environment variables which will used to configure our application and save across all the sub-commands.
	// rootCmd.PersistentFlags().StringVar(&databaseHost, "dbHost", os.Getenv("WORKERY_DB_HOST"), "The address of database.")
	// rootCmd.PersistentFlags().StringVar(&databasePort, "dbPort", os.Getenv("WORKERY_DB_PORT"), "The port of database.")
	// rootCmd.PersistentFlags().StringVar(&databaseUser, "dbUser", os.Getenv("WORKERY_DB_USER"), "The database user.")
	// rootCmd.PersistentFlags().StringVar(&databasePassword, "dbPassword", os.Getenv("WORKERY_DB_PASSWORD"), "The database password.")
	// rootCmd.PersistentFlags().StringVar(&databaseName, "dbName", os.Getenv("WORKERY_DB_NAME"), "The database name.")
	// rootCmd.PersistentFlags().StringVar(&databaseSSLMode, "databaseSSLMode", os.Getenv("WORKERY_DB_SSL_MODE"), "The database ssl mode.")
	// rootCmd.PersistentFlags().StringVar(&applicationIP, "applicationIP", os.Getenv("WORKERY_APP_IP"), "The ip address to bind this server to.")
	// rootCmd.PersistentFlags().StringVar(&applicationPort, "applicationPort", os.Getenv("WORKERY_APP_PORT"), "The port to bind this server to.")
	// rootCmd.PersistentFlags().StringVar(&applicationSigningKey, "appSignKey", os.Getenv("WORKERY_APP_SIGNING_KEY"), "The signing key.")
	// rootCmd.PersistentFlags().StringVar(&invoicebuilderPdfTemplateFilePath, "invoicebuilderPdfTemplateFilePath", os.Getenv("WORKERY_INVOICEBUILDER_PDF_TEMPLATE_FILE_PATH"), "The invoice builder pdf file location")
	// rootCmd.PersistentFlags().StringVar(&invoicebuilderData, "invoicebuilderData", os.Getenv("WORKERY_INVOICEBUILDER_DATA"), "The invoice builder save location")
	//
	// // Set the auto-migration code.
	// var b bool = false
	// if os.Getenv("WORKERY_APP_HAS_AUTO_MIGRATIONS") == "true" {
	// 	b = true
	// 	hasAutoMigrations = true
	// }
	// rootCmd.PersistentFlags().BoolVar(&hasAutoMigrations, "hasAutoMigrations", b, "The value which dictates whether to run database migrations or not.")
}

var rootCmd = &cobra.Command{
	Use:   "workery-cli",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
