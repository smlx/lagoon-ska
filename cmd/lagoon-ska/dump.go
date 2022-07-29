package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/smlx/lagoon-ska/internal/keycloak"
	"go.uber.org/zap"
)

// DumpCmd represents the `dump` command.
type DumpCmd struct {
	KeycloakAdminPassword string `kong:"env='KEYCLOAK_ADMIN_PASSWORD'"`
	KeycloakAdminUser     string `kong:"env='KEYCLOAK_ADMIN_USER'"`
	KeycloakURL           string `kong:"env='KEYCLOAK_URL'"`
}

// Run the Dump command.
func (cmd *DumpCmd) Run(log *zap.Logger) error {
	// get main process context, to cancel on SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()
	// init the keycloak client
	c, err := keycloak.NewClient(ctx, log, cmd.KeycloakURL,
		cmd.KeycloakAdminUser, cmd.KeycloakAdminPassword)
	if err != nil {
		return fmt.Errorf("couldn't init keycloak client: %v", err)
	}
	// get the raw groups JSON representation
	groups, err := c.RawGroups(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get groups: %v", err)
	}
	_, err = fmt.Println(string(groups))
	return err
}
