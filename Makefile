
version=$(shell head -n 1 VERSION)

trojan-box:
	go build -o build/trojan-box .


build-docker:
	docker build -t aresprotocollab/trojan-box:latest -f docker/Dockerfile .


