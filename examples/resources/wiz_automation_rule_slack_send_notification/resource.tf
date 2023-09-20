resource "wiz_integration_slackbot" "default" {
  name        = "default"
  slack_token = var.slack_token
  scope       = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_slackbot_send_notification" "default" {
  name           = "example"
  description    = "example description"
  enabled        = true
  integration_id = wiz_integration_slackbot.default.id
  trigger_source = "ISSUES"
  trigger_type = [
    "RESOLVED",
  ]
  filters = jsonencode({
    "severity" : [
      "CRITICAL"
    ]
  })
  slack_channel = var.slack_channel
  slack_note    = var.slack_note
}
