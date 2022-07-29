package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

// MissingGLPIDs returns true if the group has a lagoon-projects attribute, but
// is missing a group-lagoon-project-ids attribute.
func (g *Group) MissingGLPIDs() bool {
	_, lpOk := g.Attributes["lagoon-projects"]
	_, glpidOk := g.Attributes["group-lagoon-project-ids"]
	return lpOk && !glpidOk
}

// UpdateGroupLagoonProjectIDs sets the group-lagoon-project-ids attribute by
// copying the value of the lagoon-project attribute. The value of
// group-lagoon-project-ids is a JSON-encoded string.
func (g *Group) UpdateGroupLagoonProjectIDs() error {
	pids, err := g.lagoonProjectIDs()
	if err != nil {
		return fmt.Errorf("couldn't get lagoon-projects: %v", err)
	}
	value, err := json.Marshal(map[string][]int{g.Name: pids})
	if err != nil {
		return fmt.Errorf("couldn't marshal group-lagoon-project-ids")
	}
	g.Attributes["group-lagoon-project-ids"] = []string{string(value)}
	return nil
}

// lagoonProjectIDs decodes the lagoon-projects attribute and returns a slice of
// project IDs.
func (g *Group) lagoonProjectIDs() ([]int, error) {
	var pids []int
	if value, ok := g.Attributes["lagoon-projects"]; ok && len(value) > 0 {
		err := json.Unmarshal([]byte(fmt.Sprintf(`[%s]`, value[0])), &pids)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid lagoon-projects encoded attribute (%s): %v", value, err)
		}
	}
	return pids, nil
}

// Group represents a Keycloak Group. It holds the fields required when getting
// a list of groups from keycloak.
type Group struct {
	ID string `json:"id"`
	GroupUpdateRepresentation
}

// GroupUpdateRepresentation holds the fields required when updating a group.
type GroupUpdateRepresentation struct {
	Name       string              `json:"name"`
	Attributes map[string][]string `json:"attributes"`
}

// RawGroups returns the raw JSON group representation from the Keycloak API.
func (c *Client) RawGroups(ctx context.Context) ([]byte, error) {
	groupsURL := *c.baseURL
	groupsURL.Path = path.Join(c.baseURL.Path,
		"/auth/admin/realms/lagoon/groups")
	req, err := http.NewRequestWithContext(ctx, "GET", groupsURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't construct groups request: %v", err)
	}
	q := req.URL.Query()
	q.Add("briefRepresentation", "false")
	req.URL.RawQuery = q.Encode()
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't get groups: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("bad groups response: %d\n%s", res.StatusCode, body)
	}
	return io.ReadAll(res.Body)
}

// Groups returns all Keycloak Groups including their attributes.
func (c *Client) Groups(ctx context.Context) ([]Group, error) {
	rawGroups, err := c.RawGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get groups from Keycloak API: %v", err)
	}
	var groups []Group
	return groups, json.Unmarshal(rawGroups, &groups)
}
