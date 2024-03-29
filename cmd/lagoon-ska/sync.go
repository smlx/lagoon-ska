package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/smlx/lagoon-ska/internal/keycloak"
	"go.uber.org/zap"
)

// SyncCmd represents the `sync` command.
type SyncCmd struct {
	DryRun                bool   `kong:"env='DRY_RUN',default='true',help='Do not make any changes to Keycloak (default true)'"`
	KeycloakAdminPassword string `kong:"env='KEYCLOAK_ADMIN_PASSWORD'"`
	KeycloakAdminUser     string `kong:"env='KEYCLOAK_ADMIN_USER'"`
	KeycloakURL           string `kong:"env='KEYCLOAK_URL'"`
}

// Run the Sync command.
func (cmd *SyncCmd) Run(log *zap.Logger) error {
	// get main process context, to cancel on SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()
	// init the keycloak client
	c, err := keycloak.NewClient(ctx, log, cmd.KeycloakURL,
		cmd.KeycloakAdminUser, cmd.KeycloakAdminPassword)
	if err != nil {
		return fmt.Errorf("couldn't init keycloak client: %v", err)
	}
	// run the sync
	if cmd.DryRun {
		log.Info("not making any changes in dry-run mode")
	}
	return c.Sync(ctx, cmd.DryRun)
}
