AWS_ACCOUNT ?= DANGER-security

.PHONY: plan
plan:
	aws-okta exec ${AWS_ACCOUNT} -- terraform plan

.PHONY: validate 
validate:
	aws-okta exec ${AWS_ACCOUNT} -- terraform validate 

.PHONY: apply
apply:
	aws-okta exec ${AWS_ACCOUNT} -- terraform apply

.PHONY: output
output:
	aws-okta exec ${AWS_ACCOUNT} -- terraform output
