package cmd

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/mongodb"
	user_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

func init() {
	rootCmd.AddCommand(listusersCmd)
}

var listusersCmd = &cobra.Command{
	Use:   "listusers",
	Short: "Print the users",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		defaultLogger := slog.Default()
		mc := mongodb.NewStorage(cfg, defaultLogger)
		userStorer := user_ds.NewDatastore(cfg, defaultLogger, mc)

		f := &user_ds.UserListFilter{
			PageSize:  1_000_000,
			SortField: "_id",
			SortOrder: 1,
		}
		list, err := userStorer.ListByFilter(context.Background(), f)
		if err != nil {
			log.Fatalf("failed listing users with error: %v", err)
		}
		for _, user := range list.Results {
			log.Println(user.Email)
		}
	},
}

// Auto-generated comment for change 11
