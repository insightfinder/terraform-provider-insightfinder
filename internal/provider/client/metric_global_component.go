// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// GetProjectGlobalComponentID returns the project-wide component id ("Global_<projectKey>")
// used by metric escalation / alert settings to mean "all components".
//
// The projectKey is read from /api/external/v1/systemframework (needDetail=false keeps the
// response small). If the project is not listed yet (e.g. immediately after creation), it
// falls back to ComputeProjectKey, which reproduces the server's derivation.
func (c *Client) GetProjectGlobalComponentID(projectName string) (string, error) {
	key, err := c.lookupProjectKey(projectName)
	if err != nil {
		return "", err
	}
	if key == "" {
		key = ComputeProjectKey(projectName, c.Username)
	}
	return "Global_" + key, nil
}

// ComputeProjectKey derives a projectKey the same way the server does:
// hex(sha1(projectName + ownerUserName)) with leading zeros dropped.
func ComputeProjectKey(projectName, owner string) string {
	sum := sha1.Sum([]byte(projectName + owner))
	return strings.TrimLeft(hex.EncodeToString(sum[:]), "0")
}

// projectKeyCache holds projectName → projectKey from the most recent system framework
// fetch. Concurrent lookups share a single in-flight fetch, so a run touching many projects
// makes one call; a miss (e.g. a project created during this run) triggers one refetch,
// which concurrent misses also share.
type projectKeyCache struct {
	mu       sync.Mutex
	keys     map[string]string // nil until the first successful fetch
	inflight chan struct{}     // non-nil while a fetch is running; closed when it finishes
	started  int               // number of fetches started
	lastErr  error             // error of the most recent fetch
}

// lookupProjectKey returns the projectKey of projectName, or "" if the system framework
// does not list it. Results are served from the cache when present.
func (c *Client) lookupProjectKey(projectName string) (string, error) {
	pc := &c.projectKeys
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if key, ok := pc.keys[projectName]; ok {
		return key, nil
	}
	missAt := pc.started
	for {
		if pc.inflight != nil {
			// Join the running fetch instead of starting another one.
			ch, fetchNo := pc.inflight, pc.started
			pc.mu.Unlock()
			<-ch
			pc.mu.Lock()
			if key, ok := pc.keys[projectName]; ok {
				return key, nil
			}
			if fetchNo > missAt {
				// That fetch started after our miss, so its result is current.
				return "", pc.lastErr
			}
			continue // it started before our miss; fetch again
		}

		pc.started++
		ch := make(chan struct{})
		pc.inflight = ch
		pc.mu.Unlock()
		keys, err := c.fetchProjectKeys()
		pc.mu.Lock()
		if err == nil {
			pc.keys = keys
		}
		pc.lastErr = err
		pc.inflight = nil
		close(ch)
		if err != nil {
			return "", err
		}
		return pc.keys[projectName], nil
	}
}

// fetchProjectKeys returns projectName → projectKey for every project in the system
// framework. Projects owned by the provider user win over same-named shared projects.
func (c *Client) fetchProjectKeys() (map[string]string, error) {
	sf, err := c.GetSystemFramework(c.Username, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get system framework: %w", err)
	}
	keys := make(map[string]string)
	if sf == nil {
		return keys, nil
	}

	type projectDetail struct {
		ProjectName string `json:"projectName"`
		ProjectKey  string `json:"projectKey"`
		UserName    string `json:"userName"`
	}

	owned := make(map[string]bool)
	systems := append(append([]string{}, sf.OwnSystemArr...), sf.ShareSystemArr...)
	for _, systemStr := range systems {
		var system struct {
			ProjectDetailsList json.RawMessage `json:"projectDetailsList"`
			ProjectDetailList  json.RawMessage `json:"projectDetailList"`
		}
		if err := json.Unmarshal([]byte(systemStr), &system); err != nil {
			continue
		}
		raw := system.ProjectDetailsList
		if len(raw) == 0 {
			raw = system.ProjectDetailList
		}
		// The list is sometimes a JSON-encoded string rather than an array.
		var encoded string
		if err := json.Unmarshal(raw, &encoded); err == nil {
			raw = json.RawMessage(encoded)
		}
		var projects []projectDetail
		if err := json.Unmarshal(raw, &projects); err != nil {
			continue
		}
		for _, p := range projects {
			if p.ProjectName == "" || p.ProjectKey == "" || owned[p.ProjectName] {
				continue
			}
			isOwned := p.UserName == "" || p.UserName == c.Username
			if _, seen := keys[p.ProjectName]; !seen || isOwned {
				keys[p.ProjectName] = p.ProjectKey
				owned[p.ProjectName] = isOwned
			}
		}
	}
	return keys, nil
}
