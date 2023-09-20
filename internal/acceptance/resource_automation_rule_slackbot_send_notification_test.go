package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizAutomationRuleSlackBotSendNotification_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcSlackBot)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizAutomationRuleSlackBotSendNotificationBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_slackbot.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_slackbot_send_notification.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_slackbot_send_notification.foo",
						"description",
						"Provider Acceptance Test",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_slackbot_send_notification.foo",
						"enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_slackbot_send_notification.foo",
						"trigger_source",
						"ISSUES",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_slackbot_send_notification.foo",
						"trigger_type.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttr(
						"wiz_automation_rule_slackbot_send_notification.foo",
						"trigger_type.*",
						"RESOLVED",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_slackbot.foo",
						"id",
						"wiz_automation_rule_slackbot_send_notification.foo",
						"integration_id",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_slackbot_send_notification.foo",
						"slack_channel",
						"foo",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_slackbot_send_notification.foo",
						"slack_note",
						"bar",
					),
				),
			},
		},
	})
}

func testResourceWizAutomationRuleSlackBotSendNotificationBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_integration_slackbot" "foo" {
  name                = "%s"
  scope               = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_slackbot_send_notification" "foo" {
  name           = "%s"
  description    = "Provider Acceptance Test"
  enabled        = false
  integration_id = wiz_integration_slackbot.foo.id
  trigger_source = "ISSUES"
  trigger_type = [
    "RESOLVED",
  ]
  filters = jsonencode({
      "severity": [
        "CRITICAL"
      ]
    })
  slack_channel  = "foo"
  slack_note  = "bar"
}
`, rName, rName)
}
