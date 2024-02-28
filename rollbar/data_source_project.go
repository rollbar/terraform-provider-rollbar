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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Human readable name for the project",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Computed values
			"id": {
				Description: "ID of project",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"account_id": {
				Description: "ID of account that owns the project",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"date_created": {
				Description: "Date the project was created",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"date_modified": {
				Description: "Date the project was last modified",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"status": {
				Description: "Status of the project",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	c := meta.(map[string]*client.RollbarAPIClient)[schemaKeyToken]
	c.SetHeaderDataSource(rollbarProject)
 
	pl, err := c.ListProjects()
	if err != nil {
		return err
	}

	var project client.Project
	var found bool
	for _, p := range pl {
		if p.Name == name {
			found = true
			project = p
		}
	}
	if !found {
		d.SetId("")
		return fmt.Errorf("no project with the name %s found", name)
	}

	id := fmt.Sprintf("%d", project.ID)
	d.SetId(id)
	mustSet(d, "account_id", project.AccountID)
	mustSet(d, "date_created", project.DateCreated)
	mustSet(d, "date_modified", project.DateModified)
	mustSet(d, "status", project.Status)
	return nil
}
