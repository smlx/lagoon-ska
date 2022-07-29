package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"go.uber.org/zap"
)

// filterMissingAttribute returns the groups missing the
// group-lagoon-project-ids attribute
func filterMissingAttribute(groups []Group) []Group {
	var missing []Group
	for _, g := range groups {
		if g.MissingGLPIDs() {
			missing = append(missing, g)
		}
	}
	return missing
}

func constructGroupUpdateRequestBody(g *Group) (*bytes.Reader, error) {
	if len(g.Name) == 0 {
		return nil, fmt.Errorf("can't construct request with empty group name")
	}
	reqBody, err := json.Marshal(g.GroupUpdateRepresentation)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal group update request: %v", err)
	}
	return bytes.NewReader(reqBody), nil
}

func (c *Client) addGroupAttribute(ctx context.Context, g *Group) error {
	groupsURL := *c.baseURL
	groupsURL.Path = path.Join(c.baseURL.Path,
		"/auth/admin/realms/lagoon/groups", g.ID)
	reqBody, err := constructGroupUpdateRequestBody(g)
	if err != nil {
		return fmt.Errorf("couldn't construct request body: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", groupsURL.String(),
		reqBody)
	if err != nil {
		return fmt.Errorf("couldn't construct groups request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("couldn't get groups: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("bad groups response: %d\n%s", res.StatusCode, body)
	}
	body, _ := io.ReadAll(res.Body)
	if len(body) > 0 {
		c.log.Debug("group update response", zap.ByteString("body", body))
	}
	return nil
}

// Sync will synchronise the Keycloak group-lagoon-project-ids Group
// attributes.
func (c *Client) Sync(ctx context.Context, dryRun bool) error {
	// get all groups
	groups, err := c.Groups(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get groups: %v", err)
	}
	// filter for groups missing the group-lagoon-project-ids attribute
	missing := filterMissingAttribute(groups)
	c.log.Info("groups found with missing attributes",
		zap.Int("count", len(missing)))
	// add the IDs to the groups missing them
	if dryRun {
		for _, g := range missing {
			c.log.Info("missing group-lagoon-project-ids attribute on group",
				zap.String("name", g.Name), zap.String("ID", g.ID))
		}
		return nil
	}
	for _, g := range missing {
		if err = g.UpdateGroupLagoonProjectIDs(); err != nil {
			return fmt.Errorf("couldn't update group: %v", err)
		}
		err = c.addGroupAttribute(ctx, &g)
		if err != nil {
			return fmt.Errorf("couldn't update group %s (%s) in keycloak: %v",
				g.Name, g.ID, err)
		}
		c.log.Info("updated group-lagoon-project-ids attribute on group",
			zap.String("name", g.Name), zap.String("ID", g.ID))
	}
	return nil
}
