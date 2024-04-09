.PHONY: help
help:
	@echo 'make <target>'
	@echo ''
	@echo 'Targets:'
	@echo '    help                  Show this help'
	@echo '    pre-commit            Run pre-commit checks'
	@echo ''
	@echo '    gcp_example_clean     Clean up state files'
	@echo '    gcp_example_configure Configure the deployment'
	@echo '    gcp_example_deploy    Deploy configured resources'
	@echo '    gcp_example_init      Initialize modules, providers'
	@echo '    gcp_example_install   Install Terraform, Terragrunt'
	@echo '    gcp_example_lint      Run Go linters'
	@echo '    gcp_example_plan      Show deployment plan'
	@echo '    gcp_example_test      Run deployment tests'
	@echo ''

.PHONY: pre-commit
pre-commit:
	@pre-commit run -a

.PHONY: gcp_example_clean
gcp_example_clean:
	@cd gcp/example && chmod +x ./scripts/prune_caches.sh && ./scripts/prune_caches.sh .
	@cd gcp/example/test && rm -f go.mod go.sum

.PHONY: gcp_example_configure
gcp_example_configure:
	@cd gcp/example && ./scripts/configure.sh -e dev -m US -o kpeder -p us-east1 -s us-central1 -t devops

.PHONY: gcp_example_deploy
gcp_example_deploy: gcp_example_configure gcp_example_init
	@cd gcp/example/test && go test -v

.PHONY: gcp_example_init
gcp_example_init: gcp_example_configure
	@cd gcp/example && terragrunt run-all init
	@cd gcp/example/test && go mod init deployment_test.go; go mod tidy

.PHONY: gcp_example_install
gcp_example_install:
	@chmod +x ./scripts/*.sh
	@sudo ./scripts/install_terraform.sh -v ./gcp/example/versions.yaml
	@sudo ./scripts/install_terragrunt.sh -v ./gcp/example/versions.yaml

.PHONY: gcp_example_lint
gcp_example_lint: gcp_example_configure gcp_example_init
	@cd gcp/example/test && golangci-lint run --print-linter-name --verbose gcp_example_test.go

.PHONY: gcp_example_plan
gcp_example_plan: gcp_example_configure gcp_example_init
	@cd gcp/example && terragrunt run-all plan

.PHONY: gcp_example_test
gcp_example_test: gcp_example_configure gcp_example_lint gcp_example_init
	@cd gcp/example/test && go test -v -destroy
