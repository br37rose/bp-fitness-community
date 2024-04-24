package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/mongodb"
	o_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/organization/datastore"
	user_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/provider/password"
	p "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/provider/password"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	createAdminEmail    string
	createAdminPassword string
)

func init() {
	createAdminCmd.Flags().StringVarP(&createAdminEmail, "email", "e", "", "Email of the user account")
	createAdminCmd.MarkFlagRequired("email")
	createAdminCmd.Flags().StringVarP(&createAdminPassword, "password", "p", "", "Password of the user account")
	createAdminCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(createAdminCmd)
}

var createAdminCmd = &cobra.Command{
	Use:   "createadmin",
	Short: "Creates the system administrator for tenant",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		pass := password.NewProvider()
		defaultLogger := slog.Default()
		mc := mongodb.NewStorage(cfg, defaultLogger)
		orgDS := o_d.NewDatastore(cfg, defaultLogger, mc)
		userStorer := user_ds.NewDatastore(cfg, defaultLogger, mc)
		runCreateAdmin(cfg, pass, orgDS, userStorer)
	},
}

func runCreateAdmin(cfg *config.Conf, pass p.Provider, org o_d.OrganizationStorer, us user_ds.UserStorer) {
	ctx := context.Background()

	user, err := us.GetByEmail(context.Background(), changePassEmail)
	if err != nil {
		log.Fatal(err)
	}
	if user != nil {
		log.Fatal("email already exists")
	}

	tid := cfg.AppServer.InitialOrgID
	tenantID, err := primitive.ObjectIDFromHex(tid)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	t, err := org.GetByID(ctx, tenantID)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	passwordHash, err := pass.GenerateHashFromPassword(createAdminPassword)
	if err != nil {
		log.Fatal("HashPassword:", err)
	}
	m := &user_ds.User{
		OrganizationID:        t.ID,
		OrganizationName:      t.Name,
		ID:                    primitive.NewObjectID(),
		FirstName:             "Organization",
		LastName:              "Administrator",
		Name:                  "Organization Administrator",
		LexicalName:           "Administrator, Organization",
		Email:                 createAdminEmail,
		PasswordHash:          passwordHash,
		PasswordHashAlgorithm: pass.AlgorithmName(),
		Role:                  user_ds.UserRoleAdmin,
		WasEmailVerified:      true,
		CreatedAt:             time.Now(),
		ModifiedAt:            time.Now(),
		Status:                user_ds.UserStatusActive,
		AgreeTOS:              true,
	}
	err = us.Create(ctx, m)
	if err != nil {
		log.Fatal("create admin user error:", err)
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println("successfully created admin")
}
