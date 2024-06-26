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
		t.Fail()
	}
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/../local.gcp.yaml") {
		if !fileExists(terraformDeploymentOptions.TerraformDir + "/../gcp.yaml") {
			t.Fail()
		}
	}
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/reg-multi/region.yaml") {
		t.Fail()
	}
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/reg-primary/region.yaml") {
		t.Fail()
	}
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/reg-secondary/region.yaml") {
		t.Fail()
	}
	if !fileExists(terraformDeploymentOptions.TerraformDir + "/versions.yaml") {
		t.Fail()
	}

	// Define modules
	moddirs["0-exampleFolder"] = "../global/folders/example"
	moddirs["1-exampleProject"] = "../global/projects/example"
	moddirs["2-exampleAuditConfig"] = "../global/audit-configs/example"
	moddirs["2-exampleStorageBucket"] = "../reg-multi/buckets/example"

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

		fmt.Println("Validating module:", module)

		// Sanity test
		terraform.Validate(t, terraformOptions)

		// Check for standard files
		if !fileExists(terraformOptions.TerraformDir + "/inputs.yaml") {
			t.Fail()
		}
		if !fileExists(terraformOptions.TerraformDir + "/remotestate.tf") {
			t.Fail()
		}
		if !fileExists(terraformOptions.TerraformDir + "/terragrunt.hcl") {
			t.Fail()
		}
	}

	// Read and store the env.yaml
	yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/env.yaml")
	if err != nil {
		t.Fail()
	}

	env := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &env)
	if err != nil {
		t.Fail()
	}

	// Read and store the gcp.yaml
	if fileExists(terraformDeploymentOptions.TerraformDir + "/../local.gcp.yaml") {
		yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/../local.gcp.yaml")
		if err != nil {
			t.Fail()
		}
	} else {
		yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/../gcp.yaml")
		if err != nil {
			t.Fail()
		}
	}

	platform := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &platform)
	if err != nil {
		t.Fail()
	}

	// Read and store the reg-multi/region.yaml
	yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/reg-multi/region.yaml")
	if err != nil {
		t.Fail()
	}

	mregion := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &mregion)
	if err != nil {
		t.Fail()
	}

	// Read and store the reg-primary/region.yaml
	yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/reg-primary/region.yaml")
	if err != nil {
		t.Fail()
	}

	pregion := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &pregion)
	if err != nil {
		t.Fail()
	}

	// Read and store the reg-secondary/region.yaml
	yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/reg-secondary/region.yaml")
	if err != nil {
		t.Fail()
	}

	sregion := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &sregion)
	if err != nil {
		t.Fail()
	}

	// Read and store the versions.yaml
	yfile, err = os.ReadFile(terraformDeploymentOptions.TerraformDir + "/versions.yaml")
	if err != nil {
		t.Fail()
	}

	versions := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &versions)
	if err != nil {
		t.Fail()
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

		fmt.Println("Testing module:", module)

		// Read the provider output and verify configured version
		providers := terraform.RunTerraformCommand(t, terraformOptions, terraform.FormatArgs(terraformOptions, "providers")...)
		assert.Contains(t, providers, "provider[registry.terraform.io/hashicorp/google] ~> "+versions["google_provider_version"].(string))

		// Read the inputs.yaml
		yfile, err := os.ReadFile(terraformOptions.TerraformDir + "/inputs.yaml")
		if err != nil {
			t.Fail()
		}

		inputs := make(map[string]interface{})
		err = yaml.Unmarshal(yfile, &inputs)
		if err != nil {
			t.Fail()
		}

		// Read the terragrunt.hcl
		hclfile, err := os.ReadFile(terraformOptions.TerraformDir + "/terragrunt.hcl")
		if err != nil {
			t.Fail()
		}

		hclstring := string(hclfile)

		// Make sure the path referes to the correct parent hcl file
		assert.Contains(t, hclstring, "path = find_in_parent_folders(\"example_terragrunt.hcl\")")

		// Collect the outputs
		outputs := terraform.OutputAll(t, terraformOptions)

		assert.NotEmpty(t, outputs)

		// Add module-specific tests below
		// Remember that we're in a loop, so group tests by module name (modules range keys)
		// The following collections are available for tests:
		//   platform, env, mregion, pregion, sregion, versions, inputs, outputs
		// Two key patterns are available.
		// 1. Reference the output map returned by terraform.OutputAll (ie. the output of "terragrunt output")
		//		require.Equal(t, pregion["location"], outputs["location"])
		// 2. Query the json string representing state returned by terraform.Show (ie. the output of "terragrunt show -json")
		//		modulejson := gojsonq.New().JSONString(terraform.Show(t, terraformOptions)).From("values.root_module.resources").
		//			Where("address", "eq", "azurerm_resource_group.main").
		//			Select("values")
		//		// Execute the above query; since it modifies the pointer we can only do this once, so we add it to a variable
		//		values := modulejson.Get()

		// Module-specific tests
		switch module {

		// Example folder module
		case "0-exampleFolder":
			// Make sure that prevent_destroy is set to false
			assert.Contains(t, hclstring, "prevent_destroy = false")

		// Example project module
		case "1-exampleProject":
			// Make sure that prevent_destroy is set to false
			assert.Contains(t, hclstring, "prevent_destroy = false")

		// Example audit config module
		case "2-exampleAuditConfig":
			// Make sure that prevent_destroy is set to false
			assert.Contains(t, hclstring, "prevent_destroy = false")

		// Example project module
		case "2-exampleStorageBucket":
			// Make sure that prevent_destroy is set to false
			assert.Contains(t, hclstring, "prevent_destroy = false")

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
