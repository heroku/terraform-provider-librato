package librato

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heroku/go-librato/librato"
)

func resourceLibratoMetric() *schema.Resource {
	return &schema.Resource{
		Create: resourceLibratoMetricCreate,
		Read:   resourceLibratoMetricRead,
		Update: resourceLibratoMetricUpdate,
		Delete: resourceLibratoMetricDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"composite": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"color": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_max": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_min": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_units_long": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_units_short": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_stacked": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"gap_detection": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"aggregate": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceLibratoMetricCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*librato.Client)

	metric := librato.Metric{
		Name: librato.String(d.Get("name").(string)),
		Type: librato.String(d.Get("type").(string)),
	}
	if a, ok := d.GetOk("display_name"); ok {
		metric.DisplayName = librato.String(a.(string))
	}
	if a, ok := d.GetOk("description"); ok {
		metric.Description = librato.String(a.(string))
	}
	if a, ok := d.GetOk("period"); ok {
		metric.Period = librato.Uint(uint(a.(int)))
	}
	if a, ok := d.GetOk("composite"); ok {
		metric.Composite = librato.String(a.(string))
	}

	if a, ok := d.GetOk("attributes"); ok {
		attributeData := a.([]interface{})

		// attributeData SHOULD be a list of length 1 if set. Don't
		// even operate on it otherwise.
		if len(attributeData) == 1 {
			attributeDataMap := attributeData[0].(map[string]interface{})
			attributes := new(librato.MetricAttributes)

			if v, ok := attributeDataMap["color"].(string); ok && v != "" {
				attributes.Color = librato.String(v)
			}
			if v, ok := attributeDataMap["display_max"].(string); ok && v != "" {
				attributes.DisplayMax = librato.String(v)
			}
			if v, ok := attributeDataMap["display_min"].(string); ok && v != "" {
				attributes.DisplayMin = librato.String(v)
			}
			if v, ok := attributeDataMap["display_units_long"].(string); ok && v != "" {
				attributes.DisplayUnitsLong = librato.String(v)
			}
			if v, ok := attributeDataMap["display_units_short"].(string); ok && v != "" {
				attributes.DisplayUnitsShort = librato.String(v)
			}
			if v, ok := extractBool(attributeDataMap, "display_stacked"); ok {
				attributes.DisplayStacked = librato.Bool(v)
			}
			if v, ok := extractBool(attributeDataMap, "gap_detection"); ok {
				attributes.GapDetection = librato.Bool(v)
			}
			if v, ok := extractBool(attributeDataMap, "aggregate"); ok {
				attributes.Aggregate = librato.Bool(v)
			}

			metric.Attributes = attributes
		}
	}

	_, err := client.Metrics.Update(&metric)
	if err != nil {
		log.Printf("[INFO] ERROR creating Metric: %s", err)
		return fmt.Errorf("Error creating Librato metric: %s", err)
	}

	retryErr := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, _, err := client.Metrics.Get(*metric.Name)
		if err != nil {
			if errResp, ok := err.(*librato.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		return fmt.Errorf("Error creating Librato metric: %s", retryErr)
	}

	d.SetId(*metric.Name)
	return resourceLibratoMetricRead(d, meta)
}

func resourceLibratoMetricRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*librato.Client)

	id := d.Id()

	log.Printf("[INFO] Reading Librato Metric: %s", id)
	metric, _, err := client.Metrics.Get(id)
	if err != nil {
		if errResp, ok := err.(*librato.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Librato Metric %s: %s", id, err)
	}

	d.Set("name", metric.Name)
	d.Set("type", metric.Type)

	if metric.Description != nil {
		d.Set("description", metric.Description)
	}

	if metric.DisplayName != nil {
		d.Set("display_name", metric.DisplayName)
	}

	if metric.Period != nil {
		d.Set("period", metric.Period)
	}

	if metric.Composite != nil {
		d.Set("composite", metric.Composite)
	}

	attributes := metricAttributesGather(d, metric.Attributes)

	// Since attributes isn't a simple terraform type (TypeList), it's best to
	// catch the error returned from the d.Set() function, and handle accordingly.
	if err := d.Set("attributes", attributes); err != nil {
		return err
	}

	return nil
}

func resourceLibratoMetricUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*librato.Client)

	id := d.Id()

	metric := new(librato.Metric)
	metric.Name = librato.String(id)

	if d.HasChange("type") {
		metric.Type = librato.String(d.Get("type").(string))
	}
	if d.HasChange("description") {
		metric.Description = librato.String(d.Get("description").(string))
	}
	if d.HasChange("display_name") {
		metric.DisplayName = librato.String(d.Get("display_name").(string))
	}
	if d.HasChange("period") {
		metric.Period = librato.Uint(uint(d.Get("period").(int)))
	}
	if d.HasChange("composite") {
		metric.Composite = librato.String(d.Get("composite").(string))
	}
	if d.HasChange("attributes") {
		attributeData := d.Get("attributes").([]interface{})
		attributes := new(librato.MetricAttributes)

		// attributeData SHOULD be a list of length 1 if set. Don't
		// even operate on it otherwise.
		if len(attributeData) == 1 {
			attributeDataMap := attributeData[0].(map[string]interface{})

			if v, ok := attributeDataMap["color"].(string); ok && v != "" {
				attributes.Color = librato.String(v)
			}
			if v, ok := attributeDataMap["display_max"].(string); ok && v != "" {
				attributes.DisplayMax = librato.String(v)
			}
			if v, ok := attributeDataMap["display_min"].(string); ok && v != "" {
				attributes.DisplayMin = librato.String(v)
			}
			if v, ok := attributeDataMap["display_units_long"].(string); ok && v != "" {
				attributes.DisplayUnitsLong = librato.String(v)
			}
			if v, ok := attributeDataMap["display_units_short"].(string); ok && v != "" {
				attributes.DisplayUnitsShort = librato.String(v)
			}
			if v, ok := extractBool(attributeDataMap, "display_stacked"); ok {
				attributes.DisplayStacked = librato.Bool(v)
			}
			if v, ok := extractBool(attributeDataMap, "gap_detection"); ok {
				attributes.GapDetection = librato.Bool(v)
			}
			if v, ok := extractBool(attributeDataMap, "aggregate"); ok {
				attributes.Aggregate = librato.Bool(v)
			}
		}

		// Continue assigning attributes, in the event we're emptying
		// it out this will guarantee that it will be emptied out.
		metric.Attributes = attributes
	}

	log.Printf("[INFO] Updating Librato metric: %v", structToString(metric))

	_, err := client.Metrics.Update(metric)
	if err != nil {
		return fmt.Errorf("Error updating Librato metric: %s", err)
	}

	log.Printf("[INFO] Updated Librato metric %s", id)

	// Wait for propagation since Librato updates are eventually consistent
	wait := resource.StateChangeConf{
		Pending:                   []string{fmt.Sprintf("%t", false)},
		Target:                    []string{fmt.Sprintf("%t", true)},
		Timeout:                   5 * time.Minute,
		MinTimeout:                2 * time.Second,
		ContinuousTargetOccurence: 5,
		Refresh: func() (interface{}, string, error) {
			log.Printf("[INFO] Checking if Librato Metric %s was updated yet", id)
			changedMetric, _, getErr := client.Metrics.Get(id)
			if getErr != nil {
				return changedMetric, "", getErr
			}
			return changedMetric, "true", nil
		},
	}

	_, err = wait.WaitForState()
	if err != nil {
		log.Printf("[INFO] ERROR - Failed updating Librato Metric %s: %s", id, err)
		return fmt.Errorf("Failed updating Librato Metric %s: %s", id, err)
	}

	return resourceLibratoMetricRead(d, meta)
}

func resourceLibratoMetricDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*librato.Client)

	id := d.Id()

	log.Printf("[INFO] Deleting Metric: %s", id)
	_, err := client.Metrics.Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting Metric: %s", err)
	}

	log.Printf("[INFO] Verifying Metric %s deleted", id)
	retryErr := resource.Retry(1*time.Minute, func() *resource.RetryError {

		log.Printf("[INFO] Getting Metric %s", id)
		_, _, err := client.Metrics.Get(id)
		if err != nil {
			if errResp, ok := err.(*librato.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
				log.Printf("[INFO] Metric %s not found, removing from state", id)
				return nil
			}
			log.Printf("[INFO] non-retryable error attempting to Get metric: %s", err)
			return resource.NonRetryableError(err)
		}

		log.Printf("[INFO] retryable error attempting to Get metric: %s", id)
		return resource.RetryableError(fmt.Errorf("metric still exists"))
	})
	if retryErr != nil {
		return fmt.Errorf("Error deleting librato metric: %s", retryErr)
	}

	return nil
}

// Flattens an attributes hash into something that flatmap.Flatten() can handle
func metricAttributesGather(d *schema.ResourceData, attributes *librato.MetricAttributes) []interface{} {
	retAttributes := make(map[string]interface{})

	if attributes != nil {
		if attributes.Color != nil {
			retAttributes["color"] = *attributes.Color
		}
		if attributes.DisplayMax != nil {
			retAttributes["display_max"] = fmt.Sprintf("%v", attributes.DisplayMax)
		}
		if attributes.DisplayMin != nil {
			retAttributes["display_min"] = fmt.Sprintf("%v", attributes.DisplayMin)
		}
		if attributes.DisplayUnitsLong != nil {
			retAttributes["display_units_long"] = *attributes.DisplayUnitsLong
		}
		if attributes.DisplayUnitsShort != nil {
			retAttributes["display_units_short"] = *attributes.DisplayUnitsShort
		}
		if attributes.DisplayStacked != nil {
			retAttributes["display_stacked"] = *attributes.DisplayStacked
		}
		if attributes.GapDetection != nil {
			retAttributes["gap_detection"] = *attributes.GapDetection
		}
		if attributes.Aggregate != nil {
			retAttributes["aggregate"] = *attributes.Aggregate
		}
	}

	if len(retAttributes) == 0 {
		return nil
	}

	return []interface{}{retAttributes}
}

func structToString(i interface{}) string {
	s, _ := json.Marshal(i)
	return string(s)
}

func extractBool(dict map[string]interface{}, key string) (v bool, ok bool) {
	raw, ok := dict[key]
	if !ok {
		return false, false
	}

	b, ok := raw.(bool)
	return b, ok
}
