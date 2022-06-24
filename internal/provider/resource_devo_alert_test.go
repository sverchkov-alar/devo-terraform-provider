package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCreateDevoAlertResource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScaffolding,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"devo_alert.test", "name", regexp.MustCompile("myAlertEach")),
				),
			},
		},
	})
}

const testAccResourceScaffolding = `
provider "devo" {
	token = "95b302a40b80894356a55df976d2e84f"
	endpoint = "https://api-us.devo.com"
}
resource "devo_alert" "test" {
	name = "myAlertEach"
	message = "Alert definition message created with alerting API"
	description = "Alert definition description created with alerting API"
	subcategory = "lib.my.ciphertechs-prod@ciphertechs.Alert"
	query_source_code = "from vpc.aws.flow \n select *"
	correlation_trigger = "each"
	priority = 5
}
`
