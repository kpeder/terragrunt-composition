## Terragrunt Deployment Example
Example Terragrunt deployment that includes Go test routines using Terratest.

### Decision Records
This repository uses architecture decision records to record design decisions about important elements of the solution.

The ADR index is available [here](./docs/decisions/index.md).

### Requirements
Tested on Go version 1.21 on Ubuntu Linux.

Uses installed packages:
```
gcloud
golangci-lint
make
pre-commit
terraform
terragrunt
```

### Configuration
1. Install the packages listed above.
1. Make a copy of the gcp/gcp.yaml file, named local.gcp.yaml, and fill in the fields with configuration values for the target platform.
1. Use gcloud to log into the platform. Terraform uses application default credentials:
    ```
    $ gcloud auth application-default login
    ```
1. It's recommended to deploy a build project and folder first, and to use this project in the configuration for additional deployments. The build project can be deployed and managed using this framework with a couple of additional steps to create a local state configuration and then to migrate state to a remote state configuration after the project and bucket are created. ALternatively, the build resources can be pre-created and then imported into the framework using the import command. These considerations are not addressed in this example.

### Deployment
Automated installation configuration, and deployment steps are managed using Makefile targets. Use ```make help``` for a list of configured targets:
```
$ make help
make <target>

Targets:
    help                  Show this help
    pre-commit            Run pre-commit checks

    gcp_example_clean     Clean up state files
    gcp_example_configure Configure the deployment
    gcp_example_deploy    Deploy configured resources
    gcp_example_init      Initialize modules, providers
    gcp_example_install   Install Terraform, Terragrunt
    gcp_example_lint      Run Go linters
    gcp_example_plan      Show deployment plan
    gcp_example_test      Run deployment tests
```

Note that additional targets can be added in order to configure multiple environments, for example to create development and production environments.
