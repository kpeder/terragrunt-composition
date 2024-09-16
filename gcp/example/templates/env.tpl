---
environment: ENVIRONMENT
labels:
  deployment: PREFIX
  environment: ENVIRONMENT
  owner: OWNER
  team: TEAM

dependencies:
  example_folder_dependency_path: "global/folders/example"
  example_folder_mock_outputs:
    id: "123456789012"
    name: "PREFIX-ENVIRONMENT-example-folder"

  example_project_dependency_path: "global/projects/example"
  example_project_mock_outputs:
    project_id: "PREFIX-ENVIRONMENT-example"
    project_name: "PREFIX-ENVIRONMENT-example"
    project_number: "123456789012"
    service_account_email: "project-service-account@PREFIX-ENVIRONMENT-example.iam.gserviceaccount.com"

  private_network_dependency_path: "global/networks/private"
  private_network_mock_outputs:
    network_id: "PREFIX-ENVIRONMENT-private"
    project_id: "PREFIX-ENVIRONMENT-example"

  primary_subnets_dependency_path: "reg-primary/subnets/private"
  primary_subnets_mock_outputs:
    subnets:
      us-east1/primary-a:
        id: "projects/PREFIX-ENVIRONMENT-example/regions/PREGION/subnetworks/primary-a"

  secondary_subnets_dependency_path: "reg-secondary/subnets/private"
  secondary_subnets_mock_outputs:
    subnets:
      us-east1/secondary-a:
        id: "projects/PREFIX-ENVIRONMENT-example/regions/SREGION/subnetworks/secondary-a"
