/*
 * Copyright (c) 2021 Rollbar, Inc.
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
)

// resourceNotification constructs a resource representing a Rollbar notification.
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
		},
	}
}

func resourceServiceLinkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	name := d.Get("name").(string)
	template := d.Get("template").(string)

	l := log.With().Str("name", name).Logger()

	l.Info().Msg("Creating rollbar_service_link resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	sl, err := c.CreateServiceLink(name, template)
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

	l := log.With().Str("name", name).Logger()

	l.Info().Msg("Creating rollbar_service_link resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	sl, err := c.UpdateServiceLink(id, name, template)

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
	id := mustGetID(d)
	l := log.With().
		Int("id", id).
		Logger()
	l.Info().Msg("Reading rollbar_service_link resource")
	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	sl, err := c.ReadServiceLink(id)
	if err == client.ErrNotFound {
		d.SetId("")
		l.Info().Msg("Notification not found - removed from state")
		return nil
	}
	if err != nil {
		l.Err(err).Msg("error reading rollbar_notification resource")
		return diag.FromErr(err)
	}

	mustSet(d, "name", sl.Name)
	mustSet(d, "template", sl.Template)
	l.Debug().Msg("Successfully read rollbar_service_link resource")
	return nil
}

func resourceServiceLinkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := mustGetID(d)
	l := log.With().Int("id", id).Logger()
	l.Info().Msg("Deleting rollbar_service_link resource")
	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	err := c.DeleteServiceLink(id)
	if err != nil {
		l.Err(err).Msg("Error deleting rollbar_service_link resource")
		return diag.FromErr(err)
	}
	l.Debug().Msg("Successfully deleted rollbar_service_link resource")
	return nil
}
