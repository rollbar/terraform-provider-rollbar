/*
 * Copyright (c) 2022 Rollbar, Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package rollbar

import (
	"context"
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

// resourceServiceLink constructs a resource representing a Rollbar service_link.
func resourceServiceLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceLinkCreate,
		UpdateContext: resourceServiceLinkUpdate,
		ReadContext:   resourceServiceLinkRead,
		DeleteContext: resourceServiceLinkDelete,

		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Description: "Name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"template": {
				Description: "Template",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project_api_key": {
				Description: "Overrides the project_api_key defined in the provider",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceServiceLinkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	name := d.Get("name").(string)
	template := d.Get("template").(string)
	project_api_key := d.Get("project_api_key").(string)

	l := log.With().Str("name", name).Logger()

	l.Info().Msg("Creating rollbar_service_link resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	if len(project_api_key) > 0 {
		c = client.NewClient(c.BaseURL, project_api_key)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarServiceLink, c)
	sl, err := c.CreateServiceLink(name, template)
	client.Mutex.Unlock()

	if err != nil {
		l.Err(err).Send()
		d.SetId("") // removing from the state
		return diag.FromErr(err)
	}
	l = l.With().Int("id", sl.ID).Logger()

	d.SetId(strconv.Itoa(sl.ID))
	l.Debug().Msg("Successfully created rollbar_service_link resource")

	return nil
}

func resourceServiceLinkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	id := mustGetID(d)
	name := d.Get("name").(string)
	template := d.Get("template").(string)
	project_api_key := d.Get("project_api_key").(string)

	l := log.With().Str("name", name).Logger()

	l.Info().Msg("Creating rollbar_service_link resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	if len(project_api_key) > 0 {
		c = client.NewClient(c.BaseURL, project_api_key)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarServiceLink, c)
	sl, err := c.UpdateServiceLink(id, name, template)
	client.Mutex.Unlock()

	if err != nil {
		l.Err(err).Send()
		d.SetId("") // removing from the state
		return diag.FromErr(err)
	}
	if sl.ID != id {
		err = errors.New("IDs are not equal")
		l.Err(err).Send()
		d.SetId("") // removing from the state
		return diag.FromErr(err)
	}
	l = l.With().Int("id", sl.ID).Logger()

	l.Debug().Msg("Successfully updated rollbar_service_link resource")
	return nil
}

func resourceServiceLinkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	project_api_key := d.Get("project_api_key").(string)
	id := mustGetID(d)
	l := log.With().
		Int("id", id).
		Logger()
	l.Info().Msg("Reading rollbar_service_link resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	if len(project_api_key) > 0 {
		c = client.NewClient(c.BaseURL, project_api_key)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarServiceLink, c)
	sl, err := c.ReadServiceLink(id)
	client.Mutex.Unlock()

	if err == client.ErrNotFound {
		d.SetId("")
		l.Info().Msg("Service Link not found - removed from state")
		return nil
	}
	if err != nil {
		l.Err(err).Msg("error reading rollbar_service_link resource")
		return diag.FromErr(err)
	}

	mustSet(d, "name", sl.Name)
	mustSet(d, "template", sl.Template)
	l.Debug().Msg("Successfully read rollbar_service_link resource")
	return nil
}

func resourceServiceLinkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	project_api_key := d.Get("project_api_key").(string)
	id := mustGetID(d)
	l := log.With().Int("id", id).Logger()
	l.Info().Msg("Deleting rollbar_service_link resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	if len(project_api_key) > 0 {
		c = client.NewClient(c.BaseURL, project_api_key)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarServiceLink, c)
	err := c.DeleteServiceLink(id)
	client.Mutex.Unlock()

	if err != nil {
		l.Err(err).Msg("Error deleting rollbar_service_link resource")
		return diag.FromErr(err)
	}
	l.Debug().Msg("Successfully deleted rollbar_service_link resource")
	return nil
}
