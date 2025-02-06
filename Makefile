APP      := terraform-module-versions

.PHONY: README.md

README.md:
	go install . && <README.template.md subst \
		EXAMPLES_MAIN_TF="$$(cat examples/main.tf)"\
		EXAMPLE_PRETTY="$$($(APP) check examples)"\
		EXAMPLE_LIST="$$($(APP) list -o json examples | jq .)"\
		EXAMPLE_LIST_PRETTY="$$($(APP) list examples)"\
		EXAMPLE_UPDATES="$$($(APP) check -o json examples | jq .)"\
		EXAMPLE_UPDATES_PRETTY="$$($(APP) check examples)"\
		EXAMPLE_UPDATES_ALL_PRETTY="$$($(APP) check -all examples)"\
		EXAMPLE_UPDATES_SINGLE="$$($(APP) check -o json -module=consul_github_https -module=consul_github_ssh examples | jq .)"\
		EXAMPLE_UPDATES_SINGLE_ALL_PRETTY="$$($(APP) check -all -module=consul_github_https -module=consul_github_ssh examples)"\
		EXAMPLE_UPDATES_SINGLE_PRETTY="$$($(APP) check -module=consul_github_https -module=consul_github_ssh examples)"\
		USAGE="$$($(APP) -h 2>&1)"\
		USAGE_LIST="$$($(APP) list -h 2>&1)"\
		USAGE_CHECK="$$($(APP) check -h 2>&1)"\
		APP="$(APP)"> README.md
