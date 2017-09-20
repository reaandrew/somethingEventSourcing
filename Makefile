.PHONY: test-cover-html
PACKAGES = $(shell find ./ -type d -not -path '*/\.*')

test-cover-html:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

.PHONY: terraform_init
terraform_init:
	terraform init -var 'key_name=terraform' -var 'public_key_path=/home/vagrant/.ssh/id_forora_rsa.pub' etc/terraform/

.PHONY: terraform_apply
terraform_apply:
	terraform apply -var 'key_name=terraform' -var 'public_key_path=/home/vagrant/.ssh/id_forora_rsa.pub' etc/terraform/
