package auth0

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yieldr/go-auth0/management"
)

func newEmail() *schema.Resource {
	return &schema.Resource{

		Create: createEmail,
		Read:   readEmail,
		Update: updateEmail,
		Delete: deleteEmail,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_from_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"credentials": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_user": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"api_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"access_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"secret_access_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"smtp_host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"smtp_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"smtp_user": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"smtp_pass": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func createEmail(d *schema.ResourceData, m interface{}) error {
	e := buildEmail(d)
	api := m.(*management.Management)
	if err := api.Email.Create(e); err != nil {
		return err
	}
	d.SetId(e.Name)
	return nil
}

func readEmail(d *schema.ResourceData, m interface{}) error {
	api := m.(*management.Management)
	e, err := api.Email.Read(management.WithFields("name", "enabled", "default_from_address", "credentials"))
	if err != nil {
		return err
	}
	d.SetId(e.Name)
	d.Set("name", e.Name)
	d.Set("enabled", e.Enabled)
	d.Set("default_from_address", e.DefaultFromAddress)

	if credentials := e.Credentials; credentials != nil {
		d.Set("credentials", []map[string]interface{}{
			{
				"api_user":          e.Credentials.APIUser,
				"api_key":           e.Credentials.APIKey,
				"access_key_id":     e.Credentials.AccessKeyID,
				"secret_access_key": e.Credentials.SecretAccessKey,
				"region":            e.Credentials.Region,
				"smtp_host":         e.Credentials.SMTPHost,
				"smtp_port":         e.Credentials.SMTPPort,
				"smtp_user":         e.Credentials.SMTPUser,
				"smtp_pass":         e.Credentials.SMTPPass,
			},
		})
	}

	return nil
}

func updateEmail(d *schema.ResourceData, m interface{}) error {
	e := buildEmail(d)
	api := m.(*management.Management)
	return api.Email.Update(e)
}

func deleteEmail(d *schema.ResourceData, m interface{}) error {
	api := m.(*management.Management)
	return api.Email.Delete()
}

func buildEmail(d *schema.ResourceData) *management.Email {

	var credentials *management.EmailCredentials

	if v, ok := d.GetOk("credentials"); ok {

		for _, item := range v.([]interface{}) {

			item := item.(map[string]interface{})

			credentials = &management.EmailCredentials{
				APIUser:         item["api_user"].(string),
				APIKey:          item["api_key"].(string),
				AccessKeyID:     item["access_key_id"].(string),
				SecretAccessKey: item["secret_access_key"].(string),
				Region:          item["region"].(string),
				SMTPHost:        item["smtp_host"].(string),
				SMTPPort:        item["smtp_port"].(int),
				SMTPUser:        item["smtp_user"].(string),
				SMTPPass:        item["smtp_pass"].(string),
			}
		}
	}

	return &management.Email{
		Name:               d.Get("name").(string),
		Enabled:            d.Get("enabled").(bool),
		DefaultFromAddress: d.Get("default_from_address").(string),
		Credentials:        credentials,
	}
}
