package exoscaleimport

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"

	egoscale "github.com/exoscale/egoscale/v2"
	egoscalev3 "github.com/exoscale/egoscale/v3"
	exocredentials "github.com/exoscale/egoscale/v3/credentials"
)

func TestRegisterTemplate(t *testing.T) {
	const (
		zone     = "ch-gva-2"
		id       = "4dc00415-877c-4567-ae3b-dc767d1cd330"
		checksum = "098f6bcd4621d373cade4e832627b4f6"
		imageURL = "https://example.com/test.qcow2"
		success  = `{"state":"success","reference":{"id":"` + id + `"}}`
		setting  = "template_application_consistent_snapshot_enabled"
	)

	for _, tc := range []struct {
		name, option, operation, wantError string
	}{
		{name: "enabled", option: "true", operation: `{"state":"pending","id":"` + id + `"}`},
		{name: "disabled", option: "false", operation: success},
		{name: "unset", operation: success},
		{name: "failed", operation: `{"state":"failure","message":"registration failed","reference":{"id":"` + id + `"}}`, wantError: "registration failed"},
		{name: "missing_reference", operation: `{"state":"success"}`, wantError: "no template ID"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := map[string]cty.Value{
				"api_key": cty.StringVal("fake-key"), "api_secret": cty.StringVal("fake-secret"),
				"api_timeout": cty.NumberIntVal(10), "image_bucket": cty.StringVal("test"),
				"template_name":           cty.StringVal("test"),
				"template_zones":          cty.ListVal([]cty.Value{cty.StringVal(zone), cty.StringVal("ch-dk-2")}),
				"template_disable_sshkey": cty.True,
				setting:                   cty.NullVal(cty.Bool),
			}
			wantBody := map[string]any{
				"boot-mode": "legacy", "checksum": checksum, "name": "test",
				"password-enabled": true, "ssh-key-enabled": false, "url": imageURL,
			}
			if tc.option != "" {
				config[setting] = cty.BoolVal(tc.option == "true")
				wantBody["application-consistent-snapshot-enabled"] = tc.option == "true"
			}
			if tc.name == "enabled" {
				config["template_boot_mode"] = cty.StringVal("uefi")
				config["template_description"] = cty.StringVal("description")
				config["template_username"] = cty.StringVal("user")
				config["template_maintainer"] = cty.StringVal("maintainer")
				config["template_version"] = cty.StringVal("version")
				config["template_build"] = cty.StringVal("build")
				wantBody["boot-mode"] = "uefi"
				wantBody["description"] = "description"
				wantBody["default-user"] = "user"
				wantBody["maintainer"] = "maintainer"
				wantBody["version"] = "version"
				wantBody["build"] = "build"
			}

			p := new(PostProcessor)
			require.Equal(t, &hcldec.AttrSpec{Name: setting, Type: cty.Bool}, p.ConfigSpec()[setting])
			require.NoError(t, p.Configure(cty.ObjectVal(config)))

			api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Method + " " + r.URL.Path {
				case "POST /v2/template":
					var body map[string]any
					assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
					assert.Equal(t, wantBody, body)
					io.WriteString(w, tc.operation)
				case "GET /v2/operation/" + id:
					io.WriteString(w, success)
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer api.Close()

			client, err := egoscalev3.NewClient(
				exocredentials.NewStaticCredentials(p.config.APIKey, p.config.APISecret),
				egoscalev3.ClientOptWithEndpoint(egoscalev3.Endpoint(api.URL+"/v2")),
			)
			require.NoError(t, err)
			p.registration = client

			exo := new(exoscaleClientMock)
			p.exo = exo
			templateID := id
			template := &egoscale.Template{ID: &templateID}
			if tc.wantError == "" {
				exo.On("GetTemplate", mock.Anything, zone, id).Return(template, nil).Once()
			}

			ui := new(packer.MockUi)
			state := new(multistep.BasicStateBag)
			state.Put("ui", ui)
			state.Put("templates", []*egoscale.Template{})
			state.Put("image_url", imageURL)
			state.Put("image_checksum", checksum)

			action := (&stepRegisterTemplate{postProcessor: p}).Run(t.Context(), state)
			exo.AssertExpectations(t)
			if tc.wantError != "" {
				require.Equal(t, multistep.ActionHalt, action)
				require.Contains(t, ui.ErrorMessage, tc.wantError)
				require.Empty(t, state.Get("templates"))
				return
			}
			require.Equal(t, multistep.ActionContinue, action)
			require.Empty(t, ui.ErrorMessage)
			require.Equal(t, []*egoscale.Template{template}, state.Get("templates"))
		})
	}
}
