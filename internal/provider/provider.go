package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oskarm93/azurepag-client-go"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("AZUREPAG_TOKEN", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"azurepag_role_definitions":         dataSourceRoleDefinitions(),
				"azurepag_role_assignment_requests": dataSourceRoleAssignmentRequests(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"azurepag_registration": resourceRegistration(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		token := d.Get("token").(string)

		if token != "" {
			userAgent := p.UserAgent("terraform-provider-azurepag", version)
			client := azurepag.NewClient(&token, &userAgent)
			return client, diags
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "API token must be specified",
			})
			return nil, diags
		}
	}
}
