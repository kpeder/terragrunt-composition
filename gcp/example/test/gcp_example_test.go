package test

import (
	"flag"
	"fmt"
	"os"

	//regexp"
	"sort"
	//"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	//"github.com/stretchr/testify/require"
	//"github.com/thedevsaddam/gojsonq/v2"
	"gopkg.in/yaml.v3"
)

// Flag to destroy the target environment after tests
var destroy = flag.Bool("destroy", false, "destroy environment after tests")

func TestTerragruntDeployment(t *testing.T) {

	// Terraform options
	binary := "terragrunt"
	rootdir := "../."
	moddirs := make(map[string]string)

	// Non-local vars to evaluate state between modules
	var network string
	var project string
	var withGPUTemplateLink string
	var withSQLTemplateLink string
	var withWinTemplateLink string

	// Reusable vars for unmarshalling YAML files
	var err error
	var yfile []byte

	// Define the deployment root
	terraformDeploymentOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir:    rootdir,
		TerraformBinary: binary,
	})

	// Check for standard global configuration files
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/env.yaml") {
		t.Errorf("Environment configuration file NOT FOUND. Expected file %s\n", terraformDeploymentOptions.TerraformDir+"/env.yaml")
	}
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/../local.gcp.yaml") {
		if !fileExists(terraformDeploymentOptions.TerraformDir + "/../gcp.yaml") {
			t.Errorf("Platform configuration file NOT FOUND. Expected file %s\n", terraformDeploymentOptions.TerraformDir+"/../[local.]gcp.yaml")
		}
	}
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/reg-multi/region.yaml") {
		t.Errorf("Region configuration file NOT FOUND. Expected file %s\n", terraformDeploymentOptions.TerraformDir+"/reg-multi/region.yaml")
	}
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/reg-primary/region.yaml") {
		t.Errorf("Region configuration file NOT FOUND. Expected file %s\n", terraformDeploymentOptions.TerraformDir+"/reg-primary/region.yaml")
	}
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/reg-secondary/region.yaml") {
		t.Errorf("Region configuration file NOT FOUND. Expected file %s\n", terraformDeploymentOptions.TerraformDir+"/reg-secondary/region.yaml")
	}
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/versions.yaml") {
		t.Errorf("Version configuration file NOT FOUND. Expected file %s\n", terraformDeploymentOptions.TerraformDir+"/versions.yaml")
	}

	// Define modules
	moddirs["0-exampleFolder"] = "../global/folders/example"
	moddirs["1-exampleProject"] = "../global/projects/example"
	moddirs["2-exampleAuditConfig"] = "../global/audit-configs/example"
	moddirs["2-exampleMetadata"] = "../global/metadata/example"
	moddirs["2-exampleStorageBucket"] = "../reg-multi/buckets/example"
	moddirs["2-privateNetwork"] = "../global/networks/private"
	moddirs["3-primaryPrivateSubnet"] = "../reg-primary/subnets/private"
	moddirs["3-secondaryPrivateSubnet"] = "../reg-secondary/subnets/private"
	moddirs["3-serviceAccountRoles"] = "../global/roles/service-accounts/example"
	moddirs["4-instanceTemplateWithGPU"] = "../reg-primary/templates/with-gpu-tpl"
	moddirs["4-instanceTemplateWithSQL"] = "../reg-primary/templates/with-sql-tpl"
	moddirs["4-instanceTemplateWithWin"] = "../reg-secondary/templates/with-win-tpl"
	moddirs["4-primaryPrivateRouter"] = "../reg-primary/routers/private"
	moddirs["4-secondaryPrivateRouter"] = "../reg-secondary/routers/private"
	moddirs["5-instanceWithGPU"] = "../reg-primary/instances/with-gpu-inst"
	moddirs["5-instanceWithSQL"] = "../reg-primary/instances/with-sql-inst"
	moddirs["5-instanceWithWin"] = "../reg-secondary/instances/with-win-inst"

	// Maps are unsorted, so sort the keys to process the modules in order
	modkeys := make([]string, 0, len(moddirs))
	for k := range moddirs {
		modkeys = append(modkeys, k)
	}
	sort.Strings(modkeys)

	for _, module := range modkeys {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir:    moddirs[module],
			TerraformBinary: binary,
		})

		fmt.Printf("Validating module in %s\n", moddirs[module])

		// Sanity test
		terraform.Validate(t, terraformOptions)

		// Check for standard files
		if !fileExists(terraformOptions.TerraformDir + "/inputs.yaml") {
			t.Errorf("Module inputs file NOT FOUND. Expected file %s\n", terraformOptions.TerraformDir+"/inputs.yaml")
		}
		if !fileExists(terraformOptions.TerraformDir + "/remotestate.tf") {
			t.Errorf("Remotestate configuration file NOT FOUND. Expected file %s\n", terraformOptions.TerraformDir+"/remotestate.tf")
		}
		if !fileExists(terraformOptions.TerraformDir + "/terragrunt.hcl") {
			t.Errorf("Module configuration file NOT FOUND. Expected file %s\n", terraformOptions.TerraformDir+"/terragrunt.hcl")
		}
	}

	// Read and store the env.yaml
	yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/env.yaml")
	if err != nil {
		t.Errorf("Environment configuration file NOT LOADED. Expected env.yaml to be readable.")
	}

	env := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &env)
	if err != nil {
		t.Errorf("Environment configuration file NOT LOADED. Expected env.yaml to be in YAML format.")
	}

	// Read and store the gcp.yaml
	if fileExists(terraformDeploymentOptions.TerraformDir + "/../local.gcp.yaml") {
		yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/../local.gcp.yaml")
		if err != nil {
			t.Errorf("Platform configuration file NOT LOADED. Expected local.gcp.yaml to be readable.")
		}
	} else {
		yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/../gcp.yaml")
		if err != nil {
			t.Errorf("Platform configuration file NOT LOADED. Expected gcp.yaml to be readable.")
		}
	}

	platform := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &platform)
	if err != nil {
		t.Errorf("Platform configuration file NOT LOADED. Expected [local.]gcp.yaml to be in YAML format.")
	}

	// Read and store the reg-multi/region.yaml
	yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/reg-multi/region.yaml")
	if err != nil {
		t.Errorf("Region configuration file NOT LOADED. Expected reg-multi/region.yaml to be readable.")
	}

	mregion := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &mregion)
	if err != nil {
		t.Errorf("Region configuration file NOT LOADED. Expected reg-multi/region.yaml to be in YAML format.")
	}

	// Read and store the reg-primary/region.yaml
	yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/reg-primary/region.yaml")
	if err != nil {
		t.Errorf("Region configuration file NOT LOADED. Expected reg-primary/region.yaml to be readable.")
	}

	pregion := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &pregion)
	if err != nil {
		t.Errorf("Region configuration file NOT LOADED. Expected reg-multi/region.yaml to be in YAML format.")
	}

	// Read and store the reg-secondary/region.yaml
	yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/reg-secondary/region.yaml")
	if err != nil {
		t.Errorf("Region configuration file NOT LOADED. Expected reg-secondary/region.yaml to be readable.")
	}

	sregion := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &sregion)
	if err != nil {
		t.Errorf("Region configuration file NOT LOADED. Expected reg-multi/region.yaml to be in YAML format.")
	}

	// Read and store the versions.yaml
	yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/versions.yaml")
	if err != nil {
		t.Errorf("Version configuration file NOT LOADED. Expected versions.yaml to be readable.")
	}

	versions := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &versions)
	if err != nil {
		t.Errorf("Version configuration file NOT LOADED. Expected versions.yaml to be in YAML format.")
	}

	// Clean up after ourselves if flag is set
	if *destroy {
		defer terraform.TgDestroyAll(t, terraformDeploymentOptions)
	}
	// Deploy the composition
	terraform.TgApplyAll(t, terraformDeploymentOptions)

	for _, module := range modkeys {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir:    moddirs[module],
			TerraformBinary: binary,
		})

		t.Logf("Testing module in %s\n", moddirs[module])

		// Read the provider output and verify configured version
		providers := terraform.RunTerraformCommand(t, terraformOptions, terraform.FormatArgs(terraformOptions, "providers")...)
		assert.Contains(t, providers, "provider[registry.terraform.io/hashicorp/google] ~> "+versions["google_provider_version"].(string))

		// Read the inputs.yaml
		yfile, err := os.ReadFile(terraformOptions.TerraformDir + "/inputs.yaml")
		if err != nil {
			t.Errorf("Inputs file NOT LOADED. Expected inputs.yaml to be readable.")
		}

		inputs := make(map[string]interface{})
		err = yaml.Unmarshal(yfile, &inputs)
		if err != nil {
			t.Errorf("Inputs file NOT LOADED. Expected inputs.yaml to be in YAML format.")
		}

		// Read the terragrunt.hcl
		hclfile, err := os.ReadFile(terraformOptions.TerraformDir + "/terragrunt.hcl")
		if err != nil {
			t.Errorf("Terragrunt configuration file NOT LOADED. Expected terragrunt.hcl to be readable.")
		}

		hclstring := string(hclfile)

		// Make sure the path referes to the correct parent hcl file
		if !assert.Contains(t, hclstring, "path = find_in_parent_folders(\"example_terragrunt.hcl\")") {
			t.Errorf("Terragrunt parent configuration test FAILED. Expected path to parent HCL (example_terragrunt.hcl) to be configured.")
		}

		// Collect the outputs
		outputs := terraform.OutputAll(t, terraformOptions)

		if !assert.NotEmpty(t, outputs) {
			t.Errorf("Output test FAILED. Expected terragrunt output to be not empty.")
		}

		// Add module-specific tests below
		// Remember that we're in a loop, so group tests by module name (modules range keys)
		// The following collections are available for tests:
		//   platform, env, mregion, pregion, sregion, versions, inputs, outputs
		// Two key patterns are available.
		// 1. Reference the output map returned by terraform.OutputAll (ie. the output of "terragrunt output")
		//		require.Equal(t, pregion["location"], outputs["location"])
		// 2. Query the json string representing state returned by terraform.Show (ie. the output of "terragrunt show -json")
		//		modulejson := gojsonq.New().JSONString(terraform.Show(t, terraformOptions)).From("values.root_module.resources").
		//			Where("address", "eq", "resource.this").
		//			Select("values")
		//		// Execute the above query; since it modifies the pointer we can only do this once, so we add it to a variable
		//		values := modulejson.Get()

		// Module-specific tests
		switch module {

		// Example folder module
		case "0-exampleFolder":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure that the folder count matches
			if !assert.Equal(t, len(inputs["names"].([]interface{})), len(outputs["names_list"].([]interface{}))) {
				t.Errorf("Folder count test FAILED. Expected %d, got %d.", len(inputs["names"].([]interface{})), len(outputs["names_list"].([]interface{})))
			}

			// Make sure that all folder names contain prefix, environment and a configured name
			for _, n := range inputs["names"].([]interface{}) {
				name := fmt.Sprintf("%s-%s-%s", platform["prefix"].(string), env["environment"].(string), n)
				if !assert.Contains(t, outputs["names_list"].([]interface{}), name) {
					t.Errorf("Folder name test FAILED. Expected %s, to contain %s.", outputs["names_list"].([]interface{}), name)
				}
			}

		// Example project module
		case "1-exampleProject":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure the project name contains prefix, environment and a configured name
			name := fmt.Sprintf("%s-%s-%s", platform["prefix"].(string), env["environment"].(string), inputs["project_name"].(string))
			if !assert.Equal(t, outputs["project_name"].(string), name) {
				t.Errorf("Project name test FAILED. Expected %s, got %s.", name, outputs["project_name"].(string))
			}

			// Make sure that the project has a random id assigned, if configured
			if inputs["random_project_id"].(bool) {
				if !assert.Equal(t, len(outputs["project_id"].(string)), len(name)+5) {
					t.Errorf("Project random id test FAILED. Expected project_id %s to have random characters appended.", outputs["project_id"].(string))
				}
			}

			// Make sure that the API count matches
			if !assert.Equal(t, len(inputs["activate_apis"].([]interface{})), len(outputs["enabled_apis"].([]interface{}))) {
				t.Errorf("Enabled API count test FAILED. Expected %d, got %d.", len(inputs["activate_apis"].([]interface{})), len(outputs["enabled_apis"].([]interface{})))
			}

			// Make sure that all configured APIs are enabled, actually
			for _, api := range inputs["activate_apis"].([]interface{}) {
				if !assert.Contains(t, outputs["enabled_apis"].([]interface{}), api) {
					t.Errorf("Enabled APIs test FAILED. Expected %s to contain %s.", outputs["enabled_apis"].([]interface{}), api)
				}
			}

			// Store the project id
			project = outputs["project_id"].(string)

		// Example audit config module
		case "2-exampleAuditConfig":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure that the audit log config implements the proper auditing for all services
			audit_cfg := map[string]bool{
				"ADMIN_READ": false,
				"DATA_READ":  false,
				"DATA_WRITE": false,
			}
			for _, obj := range outputs["audit_log_config"].([]interface{}) {
				if assert.Equal(t, obj.(map[string]interface{})["service"].(string), "allServices") {
					switch obj.(map[string]interface{})["log_type"].(string) {

					case "ADMIN_READ":
						audit_cfg["ADMIN_READ"] = true
					case "DATA_READ":
						audit_cfg["DATA_READ"] = true
					case "DATA_WRITE":
						audit_cfg["DATA_WRITE"] = true

					}
				}
			}
			for k, b := range audit_cfg {
				if !assert.True(t, b) {
					t.Errorf("Audit configuration test FAILED. Expected log_type: %s to be configured for allServices, got %t.", k, b)
				}
			}

		// Example metadata module
		case "2-exampleMetadata":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure that the metadata contains the proper key / value pairs
			for k, v := range inputs["metadata_items"].(map[string]interface{}) {
				if !assert.Contains(t, outputs["metadata_items"].(map[string]interface{}), k) {
					t.Errorf("Metadata configuration test FAILED. Expected metadata_items to contain key %s.", k)
				} else if !assert.Equal(t, outputs["metadata_items"].(map[string]interface{})[k].(map[string]interface{})["value"].(string), v) {
					t.Errorf("Metadata configuration test FAILED. Expected metadata key %s to be set to %s.", k, inputs["metadata_items"].(map[string]interface{})[k].(string))
				}
			}

		// Example storage bucket module
		case "2-exampleStorageBucket":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure the bucket name contains prefix, environment and a configured name
			name := fmt.Sprintf("%s-%s-%s", platform["prefix"].(string), env["environment"].(string), inputs["name"].(string))
			if !assert.Equal(t, outputs["bucket"].(map[string]interface{})["name"].(string), name) {
				t.Errorf("Bucket name test FAILED. Expected %s, got %s.", name, outputs["bucket"].(map[string]interface{})["name"].(string))
			}

			// Make sure the bucket is deployed in the correct location
			if !assert.Equal(t, outputs["bucket"].(map[string]interface{})["location"].(string), mregion["location"].(string)) {
				t.Errorf("Bucket location test FAILED. Expected %s, got %s.", mregion["location"].(string), outputs["bucket"].(map[string]interface{})["location"].(string))
			}

			// Make sure versioning is correctly configured
			if !assert.Equal(t, outputs["bucket"].(map[string]interface{})["versioning"].([]interface{})[0].(map[string]interface{})["enabled"].(bool), inputs["versioning"].(bool)) {
				t.Errorf("Bucket versioning test FAILED. Expected %t, got %t.", inputs["versioning"].(bool), outputs["bucket"].(map[string]interface{})["versioning"].([]interface{})[0].(map[string]interface{})["enabled"].(bool))
			}

		// Private network module
		case "2-privateNetwork":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure the network name contains prefix, environment and a configured name
			name := fmt.Sprintf("%s-%s-%s", platform["prefix"].(string), env["environment"].(string), inputs["name"].(string))
			if !assert.Equal(t, outputs["network_name"].(string), name) {
				t.Errorf("Network name test FAILED. Expected %s, got %s.", name, outputs["network_name"].(string))
			}

			// Make sure the network is deployed to the correct project
			if !assert.Equal(t, outputs["project_id"].(string), project) {
				t.Errorf("Parent project test FAILED. Expected %s, got %s.", project, outputs["project_id"].(string))
			}

			// Store the network id
			network = outputs["network_id"].(string)

		// Primary private subnet module
		case "3-primaryPrivateSubnet":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure the subnets are properly configured
			for _, subnet := range inputs["subnets"].([]interface{}) {
				attributes := outputs["subnets"].(map[string]interface{})[pregion["region"].(string)+"/"+subnet.(map[string]interface{})["name"].(string)]

				// They belong to the correct VPC network
				if !assert.Contains(t, attributes.(map[string]interface{})["network"].(string), network) {
					t.Errorf("Parent network test FAILED for subnet %s. Expected %s to contain %s.", subnet.(map[string]interface{})["name"].(string), attributes.(map[string]interface{})["network"].(string), network)
				}

				// They are correctly named
				if !assert.Equal(t, attributes.(map[string]interface{})["name"].(string), subnet.(map[string]interface{})["name"].(string)) {
					t.Errorf("Subnet name test FAILED. Expected %s, got %s.", subnet.(map[string]interface{})["name"].(string), attributes.(map[string]interface{})["name"].(string))
				}

				// They are correctly addressed
				if !assert.Equal(t, attributes.(map[string]interface{})["ip_cidr_range"].(string), subnet.(map[string]interface{})["range"].(string)) {
					t.Errorf("Subnet address test FAILED. Expected %s, got %s.", subnet.(map[string]interface{})["range"].(string), attributes.(map[string]interface{})["ip_cidr_range"].(string))
				}

				// They are deployed to the correct region
				if !assert.Equal(t, attributes.(map[string]interface{})["region"].(string), pregion["region"].(string)) {
					t.Errorf("Subnet region test FAILED. Expected %s, got %s.", pregion["region"].(string), attributes.(map[string]interface{})["region"].(string))
				}
			}

		// Secondary private subnet module
		case "3-secondaryPrivateSubnet":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure the subnets are properly configured
			for _, subnet := range inputs["subnets"].([]interface{}) {
				attributes := outputs["subnets"].(map[string]interface{})[sregion["region"].(string)+"/"+subnet.(map[string]interface{})["name"].(string)]

				// They belong to the correct VPC network
				if !assert.Contains(t, attributes.(map[string]interface{})["network"].(string), network) {
					t.Errorf("Parent network test FAILED for subnet %s. Expected %s to contain %s.", subnet.(map[string]interface{})["name"].(string), attributes.(map[string]interface{})["network"].(string), network)
				}

				// They are correctly named
				if !assert.Equal(t, attributes.(map[string]interface{})["name"].(string), subnet.(map[string]interface{})["name"].(string)) {
					t.Errorf("Subnet name test FAILED. Expected %s, got %s.", subnet.(map[string]interface{})["name"].(string), attributes.(map[string]interface{})["name"].(string))
				}

				// They are correctly addressed
				if !assert.Equal(t, attributes.(map[string]interface{})["ip_cidr_range"].(string), subnet.(map[string]interface{})["range"].(string)) {
					t.Errorf("Subnet address test FAILED. Expected %s, got %s.", subnet.(map[string]interface{})["range"].(string), attributes.(map[string]interface{})["ip_cidr_range"].(string))
				}

				// They are deployed to the correct region
				if !assert.Equal(t, attributes.(map[string]interface{})["region"].(string), sregion["region"].(string)) {
					t.Errorf("Subnet region test FAILED. Expected %s, got %s.", sregion["region"].(string), attributes.(map[string]interface{})["region"].(string))
				}
			}

		// Service account roles
		case "3-serviceAccountRoles":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure that the role memberships are deployed to the correct project
			if !assert.Equal(t, project, outputs["project_id"].(string)) {
				t.Errorf("Target project test FAILED. Expected %s to equal %s.", outputs["project_id"].(string), project)
			}

			// Make sure that all the role memberships are created
			for _, role := range inputs["roles"].([]interface{}) {
				if !assert.Contains(t, outputs["roles"].(map[string]interface{}), role) {
					t.Errorf("Role configuration test FAILED. Expected %v to contain %s.", outputs["roles"].(map[string]interface{}), role)
				}
			}

		// Instance template with GPU
		case "4-instanceTemplateWithGPU":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure that the template is deployed to the correct project
			if !assert.Contains(t, outputs["self_link"].(string), project) {
				t.Errorf("Parent project test FAILED. Expected %s to contain %s.", outputs["self_link"].(string), project)
			}

			// Make sure that the name prefix is correct
			if !assert.Contains(t, outputs["self_link"].(string), inputs["name_prefix"].(string)) {
				t.Errorf("Parent project test FAILED. Expected %s to contain %s.", outputs["self_link"].(string), inputs["name_prefix"].(string))
			}

			// Store the self link
			withGPUTemplateLink = outputs["self_link"].(string)

		// Instance template for SQL Server
		case "4-instanceTemplateWithSQL":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure that the template is deployed to the correct project
			if !assert.Contains(t, outputs["self_link"].(string), project) {
				t.Errorf("Parent project test FAILED. Expected %s to contain %s.", outputs["self_link"].(string), project)
			}

			// Make sure that the name prefix is correct
			if !assert.Contains(t, outputs["self_link"].(string), inputs["name_prefix"].(string)) {
				t.Errorf("Name prefix test FAILED. Expected %s to contain %s.", outputs["self_link"].(string), inputs["name_prefix"].(string))
			}

			// Make sure that the network tags are correctly set
			for _, tag := range inputs["tags"].([]interface{}) {
				if !assert.Contains(t, outputs["tags"].([]interface{}), tag.(string)) {
					t.Errorf("Network tag test FAILED. Expected %v to contain %s.", outputs["tags"].([]interface{}), tag.(string))
				}
			}

			// Store the self link
			withSQLTemplateLink = outputs["self_link"].(string)

		// Instance template for Windows Server
		case "4-instanceTemplateWithWin":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure that the template is deployed to the correct project
			if !assert.Contains(t, outputs["self_link"].(string), project) {
				t.Errorf("Parent project test FAILED. Expected %s to contain %s.", outputs["self_link"].(string), project)
			}

			// Make sure that the name prefix is correct
			if !assert.Contains(t, outputs["self_link"].(string), inputs["name_prefix"].(string)) {
				t.Errorf("Name prefix test FAILED. Expected %s to contain %s.", outputs["self_link"].(string), inputs["name_prefix"].(string))
			}

			// Make sure that the network tags are correctly set
			for _, tag := range inputs["tags"].([]interface{}) {
				if !assert.Contains(t, outputs["tags"].([]interface{}), tag.(string)) {
					t.Errorf("Network tag test FAILED. Expected %v to contain %s.", outputs["tags"].([]interface{}), tag.(string))
				}
			}

			// Store the self link
			withWinTemplateLink = outputs["self_link"].(string)

		// Primary private router module
		case "4-primaryPrivateRouter":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// It belongs to the correct VPC network
			if !assert.Contains(t, outputs["router"].(map[string]interface{})["network"].(string), network) {
				t.Errorf("Parent network test FAILED for router %s. Expected %s to contain %s.", outputs["router"].(map[string]interface{})["name"].(string), outputs["router"].(map[string]interface{})["network"].(string), network)
			}

			// It is correctly named
			if !assert.Equal(t, pregion["region"].(string)+"-"+inputs["name"].(string), outputs["router"].(map[string]interface{})["name"].(string)) {
				t.Errorf("Router name test FAILED. Expected %s, got %s.", pregion["region"].(string)+inputs["name"].(string), outputs["router"].(map[string]interface{})["name"].(string))
			}

			// It is deployed to the correct region
			if !assert.Equal(t, pregion["region"].(string), outputs["router"].(map[string]interface{})["region"].(string)) {
				t.Errorf("Router region test FAILED. Expected %s, got %s.", pregion["region"].(string), outputs["router"].(map[string]interface{})["region"].(string))
			}

			// Its NAT configuration is valid
			if inputs["nat"].(bool) {
				if !assert.NotEmpty(t, outputs["nat"].(map[string]interface{})) {
					t.Errorf("NAT configuration test FAILED. Expected NAT configuration to be applied, got %v.", outputs["nat"].(map[string]interface{}))
				}
			} else {
				if !assert.Empty(t, outputs["nat"].(map[string]interface{})) {
					t.Errorf("NAT configuration test FAILED. Expected NAT configuration to be empty, got\n %v\n", outputs["nat"].(map[string]interface{}))
				}
			}

		// Secondary private router module
		case "4-secondaryPrivateRouter":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// It belongs to the correct VPC network
			if !assert.Contains(t, outputs["router"].(map[string]interface{})["network"].(string), network) {
				t.Errorf("Parent network test FAILED for router %s. Expected %s to contain %s.", outputs["router"].(map[string]interface{})["name"].(string), outputs["router"].(map[string]interface{})["network"].(string), network)
			}

			// It is correctly named
			if !assert.Equal(t, sregion["region"].(string)+"-"+inputs["name"].(string), outputs["router"].(map[string]interface{})["name"].(string)) {
				t.Errorf("Router name test FAILED. Expected %s, got %s.", sregion["region"].(string)+inputs["name"].(string), outputs["router"].(map[string]interface{})["name"].(string))
			}

			// It is deployed to the correct region
			if !assert.Equal(t, sregion["region"].(string), outputs["router"].(map[string]interface{})["region"].(string)) {
				t.Errorf("Router region test FAILED. Expected %s, got %s.", sregion["region"].(string), outputs["router"].(map[string]interface{})["region"].(string))
			}

			// Its NAT configuration is valid
			if inputs["nat"].(bool) {
				if !assert.NotEmpty(t, outputs["nat"].(map[string]interface{})) {
					t.Errorf("NAT configuration test FAILED. Expected NAT configuration to be applied, got %v.", outputs["nat"].(map[string]interface{}))
				}
			} else {
				if !assert.Empty(t, outputs["nat"].(map[string]interface{})) {
					t.Errorf("NAT configuration test FAILED. Expected NAT configuration to be empty, got\n %v\n", outputs["nat"].(map[string]interface{}))
				}
			}

		// Instance template for SQL Server
		case "5-instanceWithGPU":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure that the instance count is correct
			if !assert.Equal(t, len(outputs["instances_details"].([]interface{})), inputs["num_instances"]) {
				t.Errorf("Instance count test FAILED. Expected %d instances to be deployed. Got %d.", inputs["num_instances"].(int), len(outputs["instances_details"].([]interface{})))
			}

			// If there's an instance deployed
			if inputs["num_instances"].(int) > 0 {

				for _, instance := range outputs["instances_details"].([]interface{}) {
					// Log the instance name
					t.Logf("Testing instance %s\n", instance.(map[string]interface{})["self_link"].(string))

					// Make sure that the instance is deployed to the correct project
					if !assert.Equal(t, instance.(map[string]interface{})["project"].(string), project) {
						t.Errorf("Parent project test FAILED. Expected %s to be equal to %s.", instance.(map[string]interface{})["project"].(string), project)
					}

					// Make sure that the instance name is correct
					if !assert.Contains(t, instance.(map[string]interface{})["self_link"].(string), inputs["name"].(string)) {
						t.Errorf("Name test FAILED. Expected %s to contain %s.", instance.(map[string]interface{})["self_link"].(string), inputs["name"].(string))
					}

					// Make sure that the instance is linked to the correct template
					if !assert.Equal(t, instance.(map[string]interface{})["source_instance_template"], withGPUTemplateLink) {
						t.Errorf("Template test FAILED. Expected %s to be equal to %s.", instance.(map[string]interface{})["source_instance_template"].(string), withGPUTemplateLink)
					}

					// Make sure that the instance is deployed to the correct region and zone
					if !assert.Equal(t, instance.(map[string]interface{})["zone"].(string), pregion["region"].(string)+"-"+pregion["zone_preference"].(string)) {
						t.Errorf("Zone test FAILED. Expected %s to be equal to %s.", instance.(map[string]interface{})["zone"].(string), pregion["region"].(string)+"-"+pregion["zone_preference"].(string))
					}

					// Make sure that the instance is properly provisioned
					if !assert.Equal(t, "SPOT", instance.(map[string]interface{})["scheduling"].([]interface{})[0].(map[string]interface{})["provisioning_model"].(string)) {
						t.Errorf("Provisioning test FAILED. Expected provisioning model of instance to be SPOT, got %s.", instance.(map[string]interface{})["scheduling"].([]interface{})[0].(map[string]interface{})["[provisioning_model]"].(string))
					}

					// Make sure that the instance is in the correct state
					if !assert.Equal(t, "RUNNING", instance.(map[string]interface{})["current_status"].(string)) {
						t.Errorf("Status test FAILED. Expected status of instance to be RUNNING, got %s.", instance.(map[string]interface{})["current_status"].(string))
					}

				}
			}

		// Instance template for SQL Server
		case "5-instanceWithSQL":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure that the instance count is correct
			if !assert.Equal(t, len(outputs["instances_details"].([]interface{})), inputs["num_instances"]) {
				t.Errorf("Instance count test FAILED. Expected %d instances to be deployed. Got %d.", inputs["num_instances"].(int), len(outputs["instances_details"].([]interface{})))
			}

			// If there's an instance deployed
			if inputs["num_instances"].(int) > 0 {

				for _, instance := range outputs["instances_details"].([]interface{}) {
					// Log the instance name
					t.Logf("Testing instance %s\n", instance.(map[string]interface{})["self_link"].(string))

					// Make sure that the instance is deployed to the correct project
					if !assert.Equal(t, instance.(map[string]interface{})["project"].(string), project) {
						t.Errorf("Parent project test FAILED. Expected %s to be equal to %s.", instance.(map[string]interface{})["project"].(string), project)
					}

					// Make sure that the instance name is correct
					if !assert.Contains(t, instance.(map[string]interface{})["self_link"].(string), inputs["name"].(string)) {
						t.Errorf("Name test FAILED. Expected %s to contain %s.", instance.(map[string]interface{})["self_link"].(string), inputs["name"].(string))
					}

					// Make sure that the instance is linked to the correct template
					if !assert.Equal(t, instance.(map[string]interface{})["source_instance_template"], withSQLTemplateLink) {
						t.Errorf("Template test FAILED. Expected %s to be equal to %s.", instance.(map[string]interface{})["source_instance_template"].(string), withSQLTemplateLink)
					}

					// Make sure that the instance is deployed to the correct region and zone
					if !assert.Equal(t, instance.(map[string]interface{})["zone"].(string), pregion["region"].(string)+"-"+pregion["zone_preference"].(string)) {
						t.Errorf("Zone test FAILED. Expected %s to be equal to %s.", instance.(map[string]interface{})["zone"].(string), pregion["region"].(string)+"-"+pregion["zone_preference"].(string))
					}

					// Make sure that the instance is properly provisioned
					if !assert.Equal(t, "SPOT", instance.(map[string]interface{})["scheduling"].([]interface{})[0].(map[string]interface{})["provisioning_model"].(string)) {
						t.Errorf("Provisioning test FAILED. Expected provisioning model of instance to be SPOT, got %s.", instance.(map[string]interface{})["scheduling"].([]interface{})[0].(map[string]interface{})["[provisioning_model]"].(string))
					}

					// Make sure that the instance is in the correct state
					if !assert.Equal(t, "RUNNING", instance.(map[string]interface{})["current_status"].(string)) {
						t.Errorf("Status test FAILED. Expected status of instance to be RUNNING, got %s.", instance.(map[string]interface{})["current_status"].(string))
					}

				}
			}

		// Instance template for Windows Server
		case "5-instanceWithWin":
			// Make sure that prevent_destroy is set to false
			if !assert.Contains(t, hclstring, "prevent_destroy = false") {
				t.Errorf("HCL content test FAILED. Expected \"prevent_destroy = false\", got %s", hclstring)
			}

			// Make sure that the instance count is correct
			if !assert.Equal(t, len(outputs["instances_details"].([]interface{})), inputs["num_instances"]) {
				t.Errorf("Instance count test FAILED. Expected %d instances to be deployed. Got %d.", inputs["num_instances"].(int), len(outputs["instances_details"].([]interface{})))
			}

			// If there's an instance deployed
			if inputs["num_instances"].(int) > 0 {

				for _, instance := range outputs["instances_details"].([]interface{}) {
					// Log the instance name
					t.Logf("Testing instance %s\n", instance.(map[string]interface{})["self_link"].(string))

					// Make sure that the instance is deployed to the correct project
					if !assert.Equal(t, instance.(map[string]interface{})["project"].(string), project) {
						t.Errorf("Parent project test FAILED. Expected %s to be equal to %s.", instance.(map[string]interface{})["project"].(string), project)
					}

					// Make sure that the instance name is correct
					if !assert.Contains(t, instance.(map[string]interface{})["self_link"].(string), inputs["name"].(string)) {
						t.Errorf("Name test FAILED. Expected %s to contain %s.", instance.(map[string]interface{})["self_link"].(string), inputs["name"].(string))
					}

					// Make sure that the instance is linked to the correct template
					if !assert.Equal(t, instance.(map[string]interface{})["source_instance_template"].(string), withWinTemplateLink) {
						t.Errorf("Template test FAILED. Expected %s to be equal to %s.", instance.(map[string]interface{})["source_instance_template"].(string), withWinTemplateLink)
					}

					// Make sure that the instance is deployed to the correct region and zone
					if !assert.Equal(t, instance.(map[string]interface{})["zone"].(string), sregion["region"].(string)+"-"+sregion["zone_preference"].(string)) {
						t.Errorf("Zone test FAILED. Expected %s to be equal to %s.", instance.(map[string]interface{})["zone"].(string), sregion["region"].(string)+"-"+sregion["zone_preference"].(string))
					}

					// Make sure that the instance is properly provisioned
					if !assert.Equal(t, "SPOT", instance.(map[string]interface{})["scheduling"].([]interface{})[0].(map[string]interface{})["provisioning_model"].(string)) {
						t.Errorf("Provisioning test FAILED. Expected provisioning model of instance to be SPOT, got %s.", instance.(map[string]interface{})["scheduling"].([]interface{})[0].(map[string]interface{})["[provisioning_model]"].(string))
					}

					// Make sure that the instance is in the correct state
					if !assert.Equal(t, "RUNNING", instance.(map[string]interface{})["current_status"].(string)) {
						t.Errorf("Status test FAILED. Expected status of instance to be RUNNING, got %s.", instance.(map[string]interface{})["current_status"].(string))
					}

				}
			}
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
