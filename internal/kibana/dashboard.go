package kibana

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/kibana"
	"github.com/elastic/terraform-provider-elasticstack/internal/models"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDashboard() *schema.Resource {
	return &schema.Resource{
		Description:   "Creates a Kibana dashboard.",
		CreateContext: resourceDashboardCreate,
		UpdateContext: resourceDashboardCreate,
		ReadContext:   resourceDashboardRead,
		DeleteContext: resourceDashboardDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "A custom ID to use or a random UUID v1 or v4 will be generated and used.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"space_id": {
				Description: "An identifier for the space. If space_id is not provided, the default space is used.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				ForceNew:    true,
			},
			"attributes": {
				Description:      "The dashboard definition, this is the value that we get by exporting the dashboard from Kibana",
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: utils.DiffJsonSuppress,
				ValidateFunc:     validation.StringIsJSON,
			},
		},
	}
}

func resourceDashboardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, diags := clients.NewApiClient(d, meta)
	if diags.HasError() {
		return diags
	}
	dashboard, diags := getDashboardFromResourceData(d)
	if diags.HasError() {
		return diags
	}

	result, diags := kibana.CreateSavedObject(client, dashboard, "dashboard")
	if diags.HasError() {
		return diags
	}

	id := &clients.CompositeId{ClusterId: dashboard.SpaceID, ResourceId: result.ID}
	d.SetId(id.String())
	return resourceDashboardRead(ctx, d, meta)
}

func resourceDashboardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, diags := clients.NewApiClient(d, meta)
	if diags.HasError() {
		return diags
	}
	compId, diags := clients.CompositeIdFromStr(d.Id())
	if diags.HasError() {
		return diags
	}
	id := compId.ResourceId
	spaceId := compId.ClusterId

	dashboard, diags := kibana.GetSavedObject(client, id, spaceId, "dashboard")
	if dashboard == nil && diags == nil {
		d.SetId("")
		return diags
	}
	if diags.HasError() {
		return diags
	}

	// set the fields
	if err := d.Set("id", dashboard.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("space_id", dashboard.SpaceID); err != nil {
		return diag.FromErr(err)
	}
	attributes, err := json.Marshal(dashboard.Attributes)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("attributes", string(attributes)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDashboardDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Error("Deleting dashboards is not supported")
}

func getDashboardFromResourceData(d *schema.ResourceData) (models.SavedObject, diag.Diagnostics) {
	var diags diag.Diagnostics
	dashboard := models.SavedObject{}
	attributesStr := d.Get("attributes")
	attributes := map[string]interface{}{}
	if err := json.NewDecoder(strings.NewReader(attributesStr.(string))).Decode(&attributes); err != nil {
		return models.SavedObject{}, diag.FromErr(err)
	}
	dashboard.Attributes = attributes
	dashboard.SpaceID = d.Get("space_id").(string)
	dashboard.ID = d.Get("id").(string)
	return dashboard, diags
}
