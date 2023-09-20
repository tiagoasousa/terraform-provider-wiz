package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizIntegrationSlackBot_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TcSlackBot) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizIntegrationSlackBotBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_slackbot.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_slackbot.foo",
						"slack_token",
						os.Getenv("WIZ_INTEGRATION_SLACKBOT_TOKEN"),
					),
				),
			},
		},
	})
}

func testResourceWizIntegrationSlackBotBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_integration_slackbot" "foo" {
  name  = "%s"
  scope = "All Resources, Restrict this Integration to global roles only"
}
`, rName)
}
