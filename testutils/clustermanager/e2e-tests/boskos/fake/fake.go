/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	"fmt"

	"knative.dev/pkg/testutils/clustermanager/e2e-tests/boskos"
	boskoscommon "sigs.k8s.io/boskos/common"
)

const fakeOwner = "fake-owner"

// BoskosClient implements boskos.Operation
type BoskosClient struct {
	resources []*boskoscommon.Resource
}

func (c *BoskosClient) getOwner(host *string) string {
	if host == nil {
		return fakeOwner
	}
	return *host
}

func (c *BoskosClient) GetResources() []*boskoscommon.Resource {
	return c.resources
}

// AcquireGKEProject fakes to be no op
func (c *BoskosClient) AcquireGKEProject(resType string) (*boskoscommon.Resource, error) {
	for _, res := range c.resources {
		if res.State == boskoscommon.Free {
			res.State = boskoscommon.Busy
			res.Owner = c.getOwner(nil)
			res.Type = resType
			return res, nil
		}
	}
	return nil, fmt.Errorf("no GKE project available")
}

// ReleaseGKEProject fakes to be no op
func (c *BoskosClient) ReleaseGKEProject(name string) error {
	owner := c.getOwner(nil)
	for _, res := range c.resources {
		if res.Name == name {
			if res.Owner == owner {
				res.Owner = ""
				res.State = boskoscommon.Free
				return nil
			}
			return fmt.Errorf("got owner: %q, expect owner: %q", res.Owner, owner)
		}
	}
	return fmt.Errorf("resource doesn't exist yet: %q", name)
}

// NewGKEProject adds Boskos resources for testing purpose
func (c *BoskosClient) NewGKEProject(name string) {
	c.resources = append(c.resources, &boskoscommon.Resource{
		Type:  boskos.GKEProjectResource,
		Name:  name,
		State: boskoscommon.Free,
	})
}
