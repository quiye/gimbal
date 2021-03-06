package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	th "github.com/gophercloud/gophercloud/testhelper"
)

// CreateClusterTemplate will create a random cluster tempalte. An error will be returned if the
// cluster-template could not be created.
func CreateClusterTemplate(t *testing.T, client *gophercloud.ServiceClient) (*clustertemplates.ClusterTemplate, error) {
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		return nil, err
	}

	name := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create cluster template: %s", name)

	boolFalse := false
	createOpts := clustertemplates.CreateOpts{
		Name:                name,
		MasterFlavorID:      "m1.small",
		Public:              &boolFalse,
		ServerType:          "vm",
		ExternalNetworkID:   choices.ExternalNetworkID,
		ImageID:             choices.MagnumImageID,
		RegistryEnabled:     &boolFalse,
		DockerStorageDriver: "devicemapper",
		COE:                 "swarm",
		FlavorID:            choices.FlavorID,
		MasterLBEnabled:     &boolFalse,
		DNSNameServer:       "8.8.8.8",
		FloatingIPEnabled:   &boolFalse,
	}

	res := clustertemplates.Create(client, createOpts)
	if res.Err != nil {
		return nil, res.Err
	}

	requestID := res.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, true, requestID != "")

	t.Logf("Cluster Template %s request ID: %s", name, requestID)

	clusterTemplate, err := res.Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created cluster template: %s", clusterTemplate.Name)

	tools.PrintResource(t, clusterTemplate)
	tools.PrintResource(t, clusterTemplate.CreatedAt)

	th.AssertEquals(t, name, clusterTemplate.Name)
	th.AssertEquals(t, choices.ExternalNetworkID, clusterTemplate.ExternalNetworkID)
	th.AssertEquals(t, choices.MagnumImageID, clusterTemplate.ImageID)

	return clusterTemplate, nil
}

/*
// DeleteClusterTemplate will delete a given cluster-template. A fatal error will occur if the
// cluster-template could not be deleted. This works best as a deferred function.
func DeleteClusterTemplate(t *testing.T, client *gophercloud.ServiceClient, id string) {
	t.Logf("Attempting to delete cluster-template: %s", id)

	err := clustertemplates.Delete(client, id).ExtractErr()
	if err != nil {
		t.Fatalf("Error deleting cluster-template %s: %s:", id, err)
	}

	t.Logf("Successfully deleted cluster-template: %s", id)

	return
}
*/
