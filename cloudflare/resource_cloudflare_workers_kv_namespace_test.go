package cloudflare

import (
	"context"
	"fmt"
	"testing"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudflareWorkersKVNamespaceCreate(t *testing.T) {
	randomNamespace := acctest.RandString(10)
	namespaceName := "cloudflare_workers_kv_namespace." + randomNamespace
	resource.Test(
		t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				testAccPreCheckOrg(t)
			},
			Providers: testAccProviders,
			//CheckDestroy: testAcctCheckCloudflareWorkersKVNamespaceDestroy,
			Steps: []resource.TestStep{
				{Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudflareWorkersKVNamespaceExists(namespaceName, &cloudflare.WorkersKVNamespace{}),
				)},
			},
		},
	)
}

func testAccCheckCloudflareWorkersKVNamespaceExists(name string, route *cloudflare.WorkersKVNamespace) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Worker Route ID is set")
		}

		namespace, err := getWorkersKVNamespaceFromAPI(rs.Primary.ID)
		if err != nil {
			return err
		}

		if namespace.ID != rs.Primary.ID {
			return fmt.Errorf("Unable to find a namespace with ID: %s", rs.Primary.ID)
		}
		return nil
	}
}

func getWorkersKVNamespaceFromAPI(namespaceID string) (cloudflare.WorkersKVNamespace, error) {
	client := testAccProvider.Meta().(*cloudflare.API)
	namespaces, err := client.ListWorkersKVNamespaces(context.Background())
	if err != nil {
		return cloudflare.WorkersKVNamespace{}, err
	}

	for _, ns := range namespaces.Result {
		if ns.ID == namespaceID {
			return ns, nil
		}
	}
	return cloudflare.WorkersKVNamespace{}, fmt.Errorf("missing expected namespace id %s", namespaceID)
}

// func testAcctCheckCloudflareWorkersKVNamespaceDestroy(s *terraform.State) error {
// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "cloudflare_workers_kv_namespace" {
// 			continue
// 		}

// 		client := testAccProvider.Meta().(*cloudflare.API)
// 		namespaceID := rs.Primary.Attributes["namespace_id"]
// 		resp, err := client.DeleteWorkersKVNamespace(context.Background(), namespaceID)
// 		return err
// 	}
// }
