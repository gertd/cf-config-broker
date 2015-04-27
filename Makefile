all:	bin/cf-config-broker
	@echo "Launching at http://localhost:3000/"
	foreman start -p 3000

bin/go-env:
	GOBIN=bin go install
