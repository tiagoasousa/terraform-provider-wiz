package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizIntegrationSlackBot() *schema.Resource {
	return &schema.Resource{
		Description: "Integrations are reusable, generic connections between Wiz and third-party platforms like Slack, Google Chat, and Jira that allow data from Wiz to be passed to your preferred tool.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Identifier for this object.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the integration.",
				Required:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Identifies the date and time when the object was created.",
				Computed:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The project this action is scoped to.",
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "All Resources, Restrict this Integration to global roles only",
				Description: fmt.Sprintf(
					"Scoping to a selected Project makes this Integration accessible only to users with global roles or Project-scoped access to the selected Project. Other users will not be able to see it, use it, or view its results. Integrations restricted to global roles cannot be seen or used by users with Project-scoped roles. \n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						internal.IntegrationScope,
					),
				),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						internal.IntegrationScope,
						false,
					),
				),
			},
			"slack_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Slack Token. (default: none, environment variable: WIZ_INTEGRATION_SLACK_TOKEN)",
				DefaultFunc: schema.EnvDefaultFunc(
					"WIZ_INTEGRATION_SLACKBOT_TOKEN",
					nil,
				),
			},
		},
		CreateContext: resourceWizIntegrationSlackBotCreate,
		ReadContext:   resourceWizIntegrationSlackBotRead,
		UpdateContext: resourceWizIntegrationSlackBotUpdate,
		DeleteContext: resourceWizIntegrationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWizIntegrationSlackBotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationSlackBotCreate called...")

	// define the graphql query
	query := `mutation CreateIntegration($input: CreateIntegrationInput!) {
	  createIntegration(
	    input: $input
	  ) {
	    integration {
	      id
	    }
	  }
	}`

	vars := &wiz.CreateIntegrationInput{}
	vars.Name = d.Get("name").(string)
	vars.Type = "SLACK_BOT"
	vars.ProjectID = d.Get("project_id").(string)
	vars.IsAccessibleToAllProjects = convertIntegrationScopeToBool(d.Get("scope").(string))
	vars.Params.SlackBot = &wiz.CreateSlackBotIntegrationParamsInput{}
	vars.Params.SlackBot.Token = d.Get("slack_token").(string)

	// process the request
	data := &CreateIntegration{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_Slack", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateIntegration.Integration.ID)

	return resourceWizIntegrationSlackBotRead(ctx, d, m)
}

func resourceWizIntegrationSlackBotRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationSlackBotRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query integration (
	  $id: ID!
	) {
	  integration(
	    id: $id
	  ) {
	    id
	    name
	    createdAt
	    updatedAt
	    project {
	      id
	    }
	    type
	    isAccessibleToAllProjects
	    usedByRules {
	      id
	    }
	    paramsType: params {
	      type: __typename
	    }
	    params {
	      ... on SlackBotIntegrationParams {
	        token
	      }
	    }
	  }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	data := &ReadIntegrationPayload{}
	params := &wiz.SlackBotIntegrationParams{}
	data.Integration.Params = params
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_Slack", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		if data.Integration.ID == "" {
			tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
			tflog.Info(ctx, "Resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	// set the resource parameters
	err := d.Set("name", data.Integration.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("created_at", data.Integration.CreatedAt)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("project_id", data.Integration.Project.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("slack_token", d.Get("slack_token").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func resourceWizIntegrationSlackBotUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationSlackBotUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateIntegration(
	  $input: UpdateIntegrationInput!
	) {
	  updateIntegration(input: $input) {
	    integration {
	      id
	    }
	  }
	}`

	// populate the graphql variables
	vars := &wiz.UpdateIntegrationInput{}
	vars.ID = d.Id()
	vars.Patch.Name = d.Get("name").(string)
	vars.Patch.Params.SlackBot = &wiz.UpdateSlackBotIntegrationParamsInput{}
	vars.Patch.Params.SlackBot.Token = d.Get("slack_token").(string)

	// process the request
	data := &UpdateIntegration{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_slack", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
