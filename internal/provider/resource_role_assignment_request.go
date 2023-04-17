package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oskarm93/azurepag-client-go"
)

func resourceRoleAssignmentRequest() *schema.Resource {
	return &schema.Resource{
		Description: "TODO",

		CreateContext: resourceRoleAssignmentRequestCreate,
		ReadContext:   resourceRoleAssignmentRequestRead,
		DeleteContext: resourceRoleAssignmentRequestDelete,

		Schema: map[string]*schema.Schema{
			"role_definition_id": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"object_id": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"subject_id": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"role_name": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"assignment_state": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceRoleAssignmentRequestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*azurepag.Client)

	objectId := d.Get("object_id").(string)
	subjectId := d.Get("subject_id").(string)
	assignmentState := d.Get("assignment_state").(string)
	roleName := d.Get("role_name").(string)

	roleDefinition, err := client.GetRoleDefinition(objectId, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CreateRoleAssignmentRequest(objectId, subjectId, roleDefinition.ID, assignmentState)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceRoleAssignmentRequestRead(ctx, d, meta)

	return diags
}

func resourceRoleAssignmentRequestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*azurepag.Client)

	objectId := d.Get("object_id").(string)
	subjectId := d.Get("subject_id").(string)
	assignmentState := d.Get("assignment_state").(string)
	roleName := d.Get("role_name").(string)

	roleDefinition, err := client.GetRoleDefinition(objectId, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	roleAssignmentRequest, err := client.GetRoleAssignmentRequest(objectId, subjectId, roleDefinition.ID, assignmentState)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("role_definition_id", roleDefinition.ID)
	d.SetId(roleAssignmentRequest.ID)

	return diags
}

func resourceRoleAssignmentRequestDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*azurepag.Client)

	objectId := d.Get("object_id").(string)
	subjectId := d.Get("subject_id").(string)
	assignmentState := d.Get("assignment_state").(string)
	roleName := d.Get("role_name").(string)

	roleDefinition, err := client.GetRoleDefinition(objectId, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteRoleAssignmentRequest(objectId, subjectId, roleDefinition.ID, assignmentState)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
