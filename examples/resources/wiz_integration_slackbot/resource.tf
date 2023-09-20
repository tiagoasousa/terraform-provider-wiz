resource "wiz_integration_slackbot" "default" {
  name        = "default"
  slack_token = var.slack_token
  scope       = "All Resources, Restrict this Integration to global roles only"
}
