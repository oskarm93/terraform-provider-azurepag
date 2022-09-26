package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oskarm93/azurepag-client-go"
)

func dataSourceRoleDefinitions() *schema.Resource {
	return &schema.Resource{
		Description: "Sample data source in the Terraform provider scaffolding.",

		ReadContext: dataSourceRoleDefinitionsRead,

		Schema: map[string]*schema.Schema{
			"object_id": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeString,
				Required:    true,
			},
			"role_definitions": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRoleDefinitionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*azurepag.Client)
	objectId := d.Get("object_id").(string)
	roleDefinitionsReply, err := client.GetRoleDefinitions(objectId)
	if err != nil {
		return diag.FromErr(err)
	}

	roleDefinitions := flattenRoleDefinitions(&roleDefinitionsReply)
	err = d.Set("role_definitions", roleDefinitions)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(objectId)

	return diags
}

func flattenRoleDefinitions(roleDefinitions *[]azurepag.RoleDefinition) []interface{} {
	if roleDefinitions != nil {
		rds := make([]interface{}, len(*roleDefinitions), len(*roleDefinitions))

		for i, roleDefinition := range *roleDefinitions {
			rd := make(map[string]interface{})

			rd["id"] = roleDefinition.ID
			rd["display_name"] = roleDefinition.DisplayName
			rds[i] = rd
		}

		return rds
	}

	return make([]interface{}, 0)
}
