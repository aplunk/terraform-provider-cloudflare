package cloudflare

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCloudflareDevicePolicyZoneCertificatesCreate(t *testing.T) {
	// Temporarily unset CLOUDFLARE_API_TOKEN if it is set as the Access
	// service does not yet support the API tokens and it results in
	// misleading state error messages.
	if os.Getenv("CLOUDFLARE_API_TOKEN") != "" {
		defer func(apiToken string) {
			os.Setenv("CLOUDFLARE_API_TOKEN", apiToken)
		}(os.Getenv("CLOUDFLARE_API_TOKEN"))
		os.Setenv("CLOUDFLARE_API_TOKEN", "")
	}

	rnd := generateRandomResourceName()
	name := fmt.Sprintf("cloudflare_device_posture_integration.%s", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccessAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudflareDevicePostureIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCloudflareDevicePolicyZoneCertificates(rnd, zoneID, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "zone_id", zoneID),
					resource.TestCheckResourceAttr(name, "enable", "true"),
				),
			},
		},
	})
}

func testCloudflareDevicePolicyZoneCertificates(rnd, zoneID string, enable bool) string {
	return fmt.Sprintf(`
resource "cloudflare_device_policy_zone_certificates" "%[1]s" {
	zone_id = "%[2]s"
	enable = "%[3]t"
}
`, rnd, zoneID, enable)
}
