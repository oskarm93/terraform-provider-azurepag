package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oskarm93/azurepag-client-go"
)

func dataSourceRoleAssignmentRequests() *schema.Resource {
	return &schema.Resource{
		Description: "Sample data source in the Terraform provider scaffolding.",

		ReadContext: dataSourceRoleAssignmentRequestsRead,

		Schema: map[string]*schema.Schema{
			"object_id": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeString,
				Required:    true,
			},
			"assignment_state": {
				Description: "Assignment state of role assignment requests to pull out: Eligible or Active",
				Type:        schema.TypeString,
				Required:    true,
			},
			"role_assignment_requests": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_definition_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"assignment_state": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRoleAssignmentRequestsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*azurepag.Client)
	objectId := d.Get("object_id").(string)
	assignmentState := d.Get("assignment_state").(string)
	roleAssignmentRequestsReply, err := client.GetRoleAssignmentRequests(objectId, assignmentState)
	if err != nil {
		return diag.FromErr(err)
	}

	roleAssignmentRequests := flattenRoleAssignmentRequests(&roleAssignmentRequestsReply)
	err = d.Set("role_assignment_requests", roleAssignmentRequests)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(objectId)

	return diags
}

func flattenRoleAssignmentRequests(roleAssignmentRequests *[]azurepag.RoleAssignmentRequest) []interface{} {
	if roleAssignmentRequests != nil {
		rars := make([]interface{}, len(*roleAssignmentRequests), len(*roleAssignmentRequests))

		for i, roleAssignment := range *roleAssignmentRequests {
			rar := make(map[string]interface{})

			rar["id"] = roleAssignment.ID
			rar["group_id"] = roleAssignment.ResourceID
			rar["role_definition_id"] = roleAssignment.RoleDefinitionID
			rar["subject_id"] = roleAssignment.SubjectID
			rar["assignment_state"] = roleAssignment.AssignmentState

			rars[i] = rar
		}

		return rars
	}

	return make([]interface{}, 0)
}
