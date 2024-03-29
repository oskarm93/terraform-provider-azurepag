package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oskarm93/azurepag-client-go"
)

func resourceRegistration() *schema.Resource {
	return &schema.Resource{
		Description: "This resource ensures that an Azure AD group is registered to use Privileged Access Group feature. The group must have assignableToRoles set to true beforehand.",

		CreateContext: resourceRegistrationCreate,
		ReadContext:   resourceRegistrationRead,
		DeleteContext: resourceRegistrationDelete,

		Schema: map[string]*schema.Schema{
			"object_id": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceRegistrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*azurepag.Client)

	objectId := d.Get("object_id").(string)

	attempts := 0
	for {
		err := client.RegisterGroup(objectId)
		if err != nil {
			err_msg := err.Error()
			if err_msg != "status: 401, body: " {
				return diag.FromErr(err)
			} else {
				if attempts >= 50 {
					return diag.FromErr(err)
				}
				time.Sleep(3 * time.Second)
				attempts++
			}
		} else {
			break
		}
	}

	_, err := client.GetRoleDefinitions(objectId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(objectId)

	return diags
}

func resourceRegistrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceRegistrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Groups cannot be unregistered once registered.",
	})
	return diags
}
