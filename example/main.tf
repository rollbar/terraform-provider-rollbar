/*
 * Copyright (c) 2020 Rollbar, Inc.
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


terraform {
  required_providers {
    rollbar = {
      source  = "github.com/rollbar/rollbar"
      version = "~> 0.2"
    }
  }
}


variable "rollbar_token" {
  type = string
}

provider "rollbar" {
  api_key = var.rollbar_token
}



resource "rollbar_project" "test" {
  name = "tf-acc-test-example"
  team_ids = [rollbar_team.test_team_0.id]
  depends_on = [rollbar_team.test_team_0]
}

resource "rollbar_project_access_token" "test_1" {
  name = "test-token-1"
  project_id = rollbar_project.test.id
  scopes = ["post_client_item"]
  depends_on = [rollbar_project.test]
  rate_limit_window_size = 60
  rate_limit_window_count = 500
}

resource "rollbar_project_access_token" "test_2" {
  name = "test-token-2"
  project_id = rollbar_project.test.id
  scopes = ["post_server_item"]
  depends_on = [rollbar_project.test]
}

resource "rollbar_team" "test_team_0" {
  name = "test-team-example"
}

resource "rollbar_user" "test_user_0" {
  email = "jason.mcvetta+tf-acc-test-rollbar-provider@gmail.com"
  team_ids = [rollbar_team.test_team_0.id]
}

data "rollbar_projects" "all" {}

data "rollbar_project_access_tokens" "test" {
  project_id = rollbar_project.test.id
  prefix = "test-"
}

data "rollbar_project" "test" {
  name = rollbar_project.test.name
  depends_on = [rollbar_project.test]
}

data "rollbar_project_access_token" "test_token_1" {
  project_id = rollbar_project.test.id
  name       = "test-token-1"
  depends_on = [rollbar_project_access_token.test_1]
}
