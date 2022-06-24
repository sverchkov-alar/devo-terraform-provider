package provider

import (
	"context"
	_ "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func devoAlertEachResource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: devoAlertCreate,
		ReadContext:   devoAlertRead,
		UpdateContext: devoAlertUpdate,
		DeleteContext: devoAlertDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "Alert name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"message": {
				Description: "A short message used to identify the alert condition. This text corresponds to the Summary field in the New alert definition window of the Devo app.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": {
				Description: "The full description of the alert condition, which corresponds to the Description field in the New alert definition window of the Devo app. ",
				Type:        schema.TypeString, Optional: true,
			},
			"subcategory": {
				Description: "This value corresponds to the Subcategory field in the New alert definition window of the Devo app. ",
				Type:        schema.TypeString,
				Required:    true,
			},
			"query_source_code": {
				Description: "Specify the LINQ query whose events will launch the type of alert defined.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"priority": {
				Description: "Enter the priority level of the alert defined. Values 1 to 10 are allowed, corresponding to the default values in the application",
				Type:        schema.TypeString,
				Required:    true,
			},
			"external_offset": {
				Description: "It is used to move the main query time range backward in time. It must be expressed in milliseconds.",
				Type:        schema.TypeInt, Optional: true,
			},
			"internal_period": {
				Description: "It is used to set the subquery time range. It must be expressed in milliseconds.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"internal_offset": {
				Description: "It is used to set the subquery time range. It must be expressed in milliseconds.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"correlation_trigger": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specify how frequently you want the system to check for events matching the conditions of your query. It must be indicated in milliseconds. The minimal value is 1 second and the maximum value is 100 days.",
			},
			"threshold": {
				Description: "Specifies how many events you want to use as a limit to trigger the alert.\n\nIt must be a positive integer number.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"absolute": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set this value to true if you want to use absolute values to calculate the deviation from the median, or false if you want to use a percentage.\n\nUsing an absolute value means that the threshold specified will be considered as the number above and below which the alert will be triggered. On the other hand, using a percentage means that the threshold specified will be considered as the percentage of the median value above and below which an alert will be triggered.",
			},
			"aggregation_column": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify an aggregation column whose values will be set against the designated threshold to trigger the alert. You can choose from any of the aggregation columns added to the query but you cannot add more than one.",
			},

			"back_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies how far in the past the search extends.",
			},
		},
	}
}

func getData(d *schema.ResourceData) Alert {

	ct := CorrelationTrigger{
		Kind: d.Get("correlation_trigger").(string),
	}

	var acc AlertCorrelationContext = AlertCorrelationContext{}
	acc.CorrelationTrigger = ct
	acc.Priority = d.Get("priority").(string)
	acc.QuerySourceCode = d.Get("query_source_code").(string)

	if ct.Kind == "each" {
		if v, ok := d.GetOk("external_offset"); ok {
			acc.ExternalOffset = v.(string)
		}
		if v, ok := d.GetOk("internal_period"); ok {
			acc.InternalPeriod = v.(string)
		}
		if v, ok := d.GetOk("internal_offset"); ok {
			acc.InternalOffset = v.(string)
		}
	}

	if ct.Kind == "several" || ct.Kind == "low" {
		if v, ok := d.GetOk("period"); ok {
			acc.Period = v.(string)
		}
		if v, ok := d.GetOk("threshold"); ok {
			acc.Threshold = v.(string)
		}
	}

	if ct.Kind == "rolling" {
		if v, ok := d.GetOk("period"); ok {
			acc.Period = v.(string)
		}
		if v, ok := d.GetOk("back_period"); ok {
			acc.BackPeriod = v.(string)
		}
	}

	if ct.Kind == "deviation" || ct.Kind == "gradient" {
		if v, ok := d.GetOk("threshold"); ok {
			acc.Threshold = v.(string)
		}
		if v, ok := d.GetOk("absolute"); ok {
			acc.Absolute = v.(string)
		}
		if v, ok := d.GetOk("aggregation_column"); ok {
			acc.AggregationColumn = v.(string)
		}
	}

	alert := Alert{}
	alert.AlertCorrelationContext = acc
	if d.Id() != "" {
		alert.Id = d.Id()
	}
	if v, ok := d.GetOk("name"); ok {
		alert.Name = v.(string)
	}
	if v, ok := d.GetOk("message"); ok {
		alert.Message = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		alert.Description = v.(string)
	}
	if v, ok := d.GetOk("subcategory"); ok {
		alert.Subcategory = v.(string)
	}
	return alert
}

func devoAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)
	alert := getData(d)
	res, err := CreateAlert(alert, client.token, client.endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(res.Id)
	return nil
}

func devoAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	token := meta.(*apiClient).token
	endpoint := meta.(*apiClient).endpoint
	alerts, err := GetAlert(token, endpoint, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(alerts) == 0 {
		d.SetId("")
		return nil
	}

	alert := alerts[0]
	d.SetId(alert.Id)
	if alert.Name != "" {
		d.Set("name", alert.Name)
	}

	if alert.Message != "" {
		d.Set("message", alert.Message)
	}

	if alert.Description != "" {
		d.Set("description", alert.Description)
	}

	//if alert.Subcategory != "" {
	//	d.Set("subcategory", alert.Subcategory)
	//}

	if alert.AlertCorrelationContext.QuerySourceCode != "" {
		d.Set("query_source_code", alert.AlertCorrelationContext.QuerySourceCode)
	}

	if alert.AlertCorrelationContext.Priority != "" {
		d.Set("priority", alert.AlertCorrelationContext.Priority)
	}

	if alert.AlertCorrelationContext.CorrelationTrigger.Kind != "" {
		d.Set("correlation_trigger", alert.AlertCorrelationContext.CorrelationTrigger.Kind)
	}

	if alert.AlertCorrelationContext.ExternalOffset != "" {
		d.Set("external_offset", alert.AlertCorrelationContext.ExternalOffset)
	}

	if alert.AlertCorrelationContext.InternalPeriod != "" {
		d.Set("internal_period", alert.AlertCorrelationContext.InternalPeriod)
	}

	if alert.AlertCorrelationContext.InternalOffset != "" {
		d.Set("internal_offset", alert.AlertCorrelationContext.InternalOffset)
	}

	if alert.AlertCorrelationContext.Period != "" {
		d.Set("period", alert.AlertCorrelationContext.Period)
	}

	if alert.AlertCorrelationContext.Threshold != "" {
		d.Set("threshold", alert.AlertCorrelationContext.Threshold)
	}

	if alert.AlertCorrelationContext.BackPeriod != "" {
		d.Set("back_period", alert.AlertCorrelationContext.BackPeriod)
	}

	if alert.AlertCorrelationContext.Absolute != "" {
		d.Set("absolute", alert.AlertCorrelationContext.Absolute)
	}
	if alert.AlertCorrelationContext.AggregationColumn != "" {
		d.Set("aggregation_column", alert.AlertCorrelationContext.AggregationColumn)
	}

	return nil
}

func devoAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)
	alert := getData(d)
	_, err := UpdateAlert(alert, client.token, client.endpoint)
	if err != nil {
		diag.FromErr(err)
	}
	return nil
}

func devoAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)
	err := DeleteAlert(client.token, client.endpoint, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
