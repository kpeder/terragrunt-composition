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

  with_gpu_template_dependency_path: "reg-primary/templates/with-gpu-tpl"
  with_gpu_template_mock_outputs:
    self_link: https://www.googleapis.com/compute/beta/projects/PREFIX-ENVIRONMENT-example/global/instanceTemplates/n1-standard-with-t4-gpu-202409161234567800000001

  with_sql_template_dependency_path: "reg-primary/templates/with-sql-tpl"
  with_sql_template_mock_outputs:
    self_link: https://www.googleapis.com/compute/beta/projects/PREFIX-ENVIRONMENT-example/global/instanceTemplates/n1-highmem-8-with-sql-202409161234567800000001

  with_gpu_instance_dependency_path: "reg-primary/instances/with-gpu-inst"
  with_gpu_instance_mock_outputs:
    instances_self_links:
      - "https://www.googleapis.com/compute/v1/projects/PREFIX-ENVIRONMENT-example/zones/PREGION-a/instances/jammy-jellyfish-001"

  with_sql_instance_dependency_path: "reg-primary/instances/with-sql-inst"
  with_sql_instance_mock_outputs:
    instances_self_links:
      - "https://www.googleapis.com/compute/v1/projects/PREFIX-ENVIRONMENT-example/zones/PREGION-a/instances/sqlsvrstd-001"
