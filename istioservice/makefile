all: unload deploy

unload:
	kubectl -n dev delete -f services.yaml

build:
	include ./service-a/makefile
	make all

	include ./service-b/makefile
	make all

deploy:
	kubectl -n dev apply -f services.yaml
