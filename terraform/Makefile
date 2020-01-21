check-env:
ifndef AWS_ACCOUNT
	$(error AWS_ACCOUNT is undefined)
endif

.PHONY: plan
plan: check-env
	aws-okta exec ${AWS_ACCOUNT} -- terraform $@

.PHONY: validate 
validate: check-env
	aws-okta exec ${AWS_ACCOUNT} -- terraform $@

.PHONY: apply
apply: check-env
	aws-okta exec ${AWS_ACCOUNT} -- terraform $@

.PHONY: destroy
destroy: check-env
	aws-okta exec ${AWS_ACCOUNT} -- terraform $@

.PHONY: output
output: check-env
	aws-okta exec ${AWS_ACCOUNT} -- terraform $@

.PHONY: fmt
fmt: check-env
	aws-okta exec ${AWS_ACCOUNT} -- terraform $@

.PHONY: graph
graph: check-env
	aws-okta exec ${AWS_ACCOUNT} -- terraform $@ | dot -Tsvg > graph.svg
