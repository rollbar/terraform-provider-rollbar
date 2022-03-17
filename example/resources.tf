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

resource "rollbar_team" "test_team_1" {
  name = "test-team-example_1"
}

resource "rollbar_team" "test_team_0" {
  name = "test-team-example"
}

resource "rollbar_user" "test_user_0" {
  email = "jason.mcvetta+tf-acc-test-rollbar-provider@gmail.com"
  team_ids = [rollbar_team.test_team_0.id]
}

resource "rollbar_team_user" "test_team_user" {
  email   = "example+tf-acc-test-rollbar-provider@gmail.com"
  team_id = rollbar_team.test_team_0.id
}

resource "rollbar_notification" "slack_notification" {
  rule  {
    filters {
        type =  "environment"
        operation =  "eq"
        value = "production"
    }
    filters {
       type = "framework"
       operation = "eq"
       value = 13
    }
   trigger = "new_item"
  }
  channel = "slack"
  config  {
     show_message_buttons = true
     channel = "#demo-user"
  }
}
resource "rollbar_service_link" "service_link" {
  name = "some_name_some_name"
  template = "sometemplate_new.{{ss}}"
}