package provider

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oskarm93/azurepag-client-go"
)

type RoleSettingsOptions struct {
	AllowPermanentEligibleAssignments bool
	MaxEligibleAssignmentTimeMins     int
	MaxActivationTimeMins             int
	RequireMFAOnActivation            bool
	RequireJustificationOnActivation  bool
	RequireTicketInfoOnActivation     bool
}

func resourceRoleSettings() *schema.Resource {
	return &schema.Resource{
		Description: "TODO",

		CreateContext: resourceRoleSettingsCreate,
		ReadContext:   resourceRoleSettingsRead,
		DeleteContext: resourceRoleSettingsDelete,
		UpdateContext: resourceRoleSettingsUpdate,

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
			"role_name": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"allow_permanent_eligible_assignments": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"max_eligible_assignment_time_mins": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"max_activation_time_mins": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"require_mfa_on_activation": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"require_justification_on_activation": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"require_ticket_info_on_activation": {
				Description: "Object ID of the Azure AD group",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceRoleSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*azurepag.Client)

	objectId := d.Get("object_id").(string)
	roleName := d.Get("role_name").(string)

	roleDefinition, err := client.GetRoleDefinition(objectId, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	existingRoleSettings, err := client.GetRoleSettings(objectId, roleDefinition.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	updatedRoleSettings, err := createUpdatedRoleSettings(existingRoleSettings, d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.UpdateRoleSettings(updatedRoleSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRoleSettingsRead(ctx, d, meta)
}

func resourceRoleSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*azurepag.Client)

	objectId := d.Get("object_id").(string)
	roleName := d.Get("role_name").(string)

	roleDefinition, err := client.GetRoleDefinition(objectId, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	roleSettings, err := client.GetRoleSettings(objectId, roleDefinition.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	roleSettingsOptions, err := getRoleSettingsOptions(roleSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("allow_permanent_eligible_assignments", roleSettingsOptions.AllowPermanentEligibleAssignments)
	d.Set("max_eligible_assignment_time_mins", roleSettingsOptions.MaxEligibleAssignmentTimeMins)
	d.Set("max_activation_time_mins", roleSettingsOptions.MaxActivationTimeMins)
	d.Set("require_justification_on_activation", roleSettingsOptions.RequireJustificationOnActivation)
	d.Set("require_mfa_on_activation", roleSettingsOptions.RequireMFAOnActivation)
	d.Set("require_ticket_info_on_activation", roleSettingsOptions.RequireTicketInfoOnActivation)
	d.Set("role_definition_id", roleDefinition.ID)
	d.SetId(roleSettings.ID)

	return diags
}

func resourceRoleSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceRoleSettingsCreate(ctx, d, meta)
}

func resourceRoleSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO: Should we restore original settings ?
	d.SetId("")

	return diags
}

func createUpdatedRoleSettings(roleSettings *azurepag.RoleSettings, d *schema.ResourceData) (*azurepag.RoleSettings, error) {
	roleSettingsOptions, err := getRoleSettingsOptions(roleSettings)
	if err != nil {
		return nil, err
	}

	assignmentExpirationRuleSetting, err := json.Marshal(azurepag.RoleSettingsExpirationRuleSetting{
		PermanentAssignment: func() bool {
			if d.HasChange("allow_permanent_eligible_assignments") {
				return d.Get("allow_permanent_eligible_assignments").(bool)
			} else {
				return roleSettingsOptions.AllowPermanentEligibleAssignments
			}
		}(),
		MaximumGrantPeriodInMinutes: func() int {
			if d.HasChange("max_eligible_assignment_time_mins") {
				return d.Get("max_eligible_assignment_time_mins").(int)
			} else {
				return roleSettingsOptions.MaxEligibleAssignmentTimeMins
			}
		}(),
	})
	if err != nil {
		return nil, err
	}

	expirationRuleSetting, err := json.Marshal(azurepag.RoleSettingsExpirationRuleSetting{
		PermanentAssignment: true, // This property is unused but must be specified
		MaximumGrantPeriodInMinutes: func() int {
			if d.HasChange("max_activation_time_mins") {
				return d.Get("max_activation_time_mins").(int)
			} else {
				return roleSettingsOptions.MaxActivationTimeMins
			}
		}(),
	})
	if err != nil {
		return nil, err
	}

	mfaRuleSetting, err := json.Marshal(azurepag.RoleSettingsMfaRuleSetting{
		MFARequired: func() bool {
			if d.HasChange("require_mfa_on_activation") {
				return d.Get("require_mfa_on_activation").(bool)
			} else {
				return roleSettingsOptions.RequireMFAOnActivation
			}
		}(),
	})
	if err != nil {
		return nil, err
	}

	justificationRuleSetting, err := json.Marshal(azurepag.RoleSettingsJustificationRuleSetting{
		Required: func() bool {
			if d.HasChange("require_justification_on_activation") {
				return d.Get("require_justification_on_activation").(bool)
			} else {
				return roleSettingsOptions.RequireJustificationOnActivation
			}
		}(),
	})
	if err != nil {
		return nil, err
	}

	ticketInfoRuleSetting, err := json.Marshal(azurepag.RoleSettingsTicketingRuleSetting{
		TicketingRequired: func() bool {
			if d.HasChange("require_ticket_info_on_activation") {
				return d.Get("require_ticket_info_on_activation").(bool)
			} else {
				return roleSettingsOptions.RequireTicketInfoOnActivation
			}
		}(),
	})
	if err != nil {
		return nil, err
	}

	result := azurepag.RoleSettings{
		ID: roleSettings.ID,
		LifecycleManagement: []azurepag.LifecycleManagement{
			{
				Caller:    "EndUser",
				Level:     "Member",
				Operation: "ALL",
				RoleSettingsRules: []azurepag.RoleSettingsRule{
					{
						RuleIdentifier: "ExpirationRule",
						Setting:        string(expirationRuleSetting),
					},
					{
						RuleIdentifier: "MfaRule",
						Setting:        string(mfaRuleSetting),
					},
					{
						RuleIdentifier: "JustificationRule",
						Setting:        string(justificationRuleSetting),
					},
					{
						RuleIdentifier: "TicketingRule",
						Setting:        string(ticketInfoRuleSetting),
					},
				},
			},
			{
				Caller:    "Admin",
				Level:     "Eligible",
				Operation: "ALL",
				RoleSettingsRules: []azurepag.RoleSettingsRule{
					{
						RuleIdentifier: "ExpirationRule",
						Setting:        string(assignmentExpirationRuleSetting),
					},
				},
			},
		},
	}

	return &result, nil
}

func getRoleSettingsOptions(roleSettings *azurepag.RoleSettings) (*RoleSettingsOptions, error) {
	activationRules, err := getActivationRules(roleSettings)
	if err != nil {
		return nil, err
	}

	eligibleAssignmentRules, err := getEligibleAssignmentRules(roleSettings)
	if err != nil {
		return nil, err
	}

	ticketingRule, err := getRuleSetting(activationRules.RoleSettingsRules, "TicketingRule")
	if err != nil {
		return nil, err
	}

	ticketingRuleSetting := azurepag.RoleSettingsTicketingRuleSetting{}
	err = json.Unmarshal([]byte(ticketingRule.Setting), &ticketingRuleSetting)
	if err != nil {
		return nil, err
	}

	mfaRule, err := getRuleSetting(activationRules.RoleSettingsRules, "MfaRule")
	if err != nil {
		return nil, err
	}

	mfaRuleSetting := azurepag.RoleSettingsMfaRuleSetting{}
	err = json.Unmarshal([]byte(mfaRule.Setting), &mfaRuleSetting)
	if err != nil {
		return nil, err
	}

	justificationRule, err := getRuleSetting(activationRules.RoleSettingsRules, "JustificationRule")
	if err != nil {
		return nil, err
	}

	justificationRuleSetting := azurepag.RoleSettingsJustificationRuleSetting{}
	err = json.Unmarshal([]byte(justificationRule.Setting), &justificationRuleSetting)
	if err != nil {
		return nil, err
	}

	expirationRule, err := getRuleSetting(activationRules.RoleSettingsRules, "ExpirationRule")
	if err != nil {
		return nil, err
	}

	expirationRuleSettings := azurepag.RoleSettingsExpirationRuleSetting{}
	err = json.Unmarshal([]byte(expirationRule.Setting), &expirationRuleSettings)
	if err != nil {
		return nil, err
	}

	assignmentExpirationRule, err := getRuleSetting(eligibleAssignmentRules.RoleSettingsRules, "ExpirationRule")
	if err != nil {
		return nil, err
	}

	assignmentExpirationRuleSettings := azurepag.RoleSettingsExpirationRuleSetting{}
	err = json.Unmarshal([]byte(assignmentExpirationRule.Setting), &assignmentExpirationRuleSettings)
	if err != nil {
		return nil, err
	}

	result := RoleSettingsOptions{
		MaxEligibleAssignmentTimeMins:     assignmentExpirationRuleSettings.MaximumGrantPeriodInMinutes,
		AllowPermanentEligibleAssignments: assignmentExpirationRuleSettings.PermanentAssignment,
		MaxActivationTimeMins:             expirationRuleSettings.MaximumGrantPeriodInMinutes,
		RequireMFAOnActivation:            mfaRuleSetting.MFARequired,
		RequireJustificationOnActivation:  justificationRuleSetting.Required,
		RequireTicketInfoOnActivation:     ticketingRuleSetting.TicketingRequired,
	}

	return &result, nil
}

func getRuleSetting(rules []azurepag.RoleSettingsRule, ruleIdentifier string) (*azurepag.RoleSettingsRule, error) {
	for _, item := range rules {
		if item.RuleIdentifier == ruleIdentifier {
			return &item, nil
		}
	}
	return nil, errors.New("Role activation rules not found.")
}

func getActivationRules(roleSettings *azurepag.RoleSettings) (*azurepag.LifecycleManagement, error) {
	for _, item := range roleSettings.LifecycleManagement {
		if item.Caller == "EndUser" && item.Level == "Member" && item.Operation == "ALL" {
			return &item, nil
		}
	}
	return nil, errors.New("Role activation rules not found.")
}

func getEligibleAssignmentRules(roleSettings *azurepag.RoleSettings) (*azurepag.LifecycleManagement, error) {
	for _, item := range roleSettings.LifecycleManagement {
		if item.Caller == "Admin" && item.Level == "Eligible" && item.Operation == "ALL" {
			return &item, nil
		}
	}
	return nil, errors.New("Role eligible assignment rules not found.")
}
