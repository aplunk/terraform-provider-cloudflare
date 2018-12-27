package cloudflare

// func TestAccCloudflareWorkersKVNamespaceCreate(t *testing.T) {
// 	resource.Test(
// 		t, resource.TestCase{
// 			PreCheck: func() {
// 				testAccPreCheck(t)
// 				testAccPreCheckOrg(t)
// 			},
// 			Providers:    testAccProviders,
// 			CheckDestroy: testAcctCheckCloudflareWorkersKVNamespaceDestroy,
// 			Steps: []resource.TestStep{
// 				test
// 			},
// 		},
// 	)
// }

// func testAccCheckCloudflareWorkersKVNamespaceExists(s *terraform.State) error {
// 	return func(s *ter)
// }

// func testAcctCheckCloudflareWorkersKVNamespaceDestroy(s *terraform.State) error {
// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "cloudflare_workers_kv_namespace" {
// 			continue
// 		}

// 		client := testAccProvider.Meta().(*cloudflare.API)
// 		namespaceID := rs.Primary.Attributes["namespace_id"]
// 		_, err := client.DeleteWorkersKVNamespace(context.Background(), namespaceID)
// 		return err
// 	}
// }
