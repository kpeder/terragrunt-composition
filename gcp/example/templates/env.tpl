---
environment: ENVIRONMENT
labels:
  deployment: PREFIX
  environment: ENVIRONMENT
  owner: OWNER
  team: TEAM
locations:
  multiregion: MREGION
  primary: PREGION
  secondary: SREGION

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
