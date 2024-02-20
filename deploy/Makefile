CUR_DIR = $(shell pwd)
TERRAFORM_DIR = $(CUR_DIR)/terraform

init:
	cd ${TERRAFORM_DIR} \
	&& terraform init

plan-dev:
	cd ${TERRAFORM_DIR} \
	&& export TF_VAR_env=dev \
	&& terraform plan

apply-dev:
	cd ${TERRAFORM_DIR} \
	&& export TF_VAR_env=dev \
	&& terraform apply -auto-approve

plan-destroy-dev:
	cd ${TERRAFORM_DIR} \
	&& export TF_VAR_env=dev \
	&& terraform plan -destroy

destroy-dev:
	cd ${TERRAFORM_DIR} \
	&& export TF_VAR_env=dev \
	&& terraform apply -destroy -auto-approve