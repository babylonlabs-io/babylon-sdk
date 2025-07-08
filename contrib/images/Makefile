BABYLON_SDK_FULL_PATH := $(shell git rev-parse --show-toplevel)

babylond: babylond-rmi
	docker build --tag babylonlabs-io/babylond -f babylond/Dockerfile ${BABYLON_SDK_FULL_PATH}

babylond-rmi:
	docker rmi babylonlabs-io/babylond --force 2>/dev/null; true

ibcsim-bcd:
	docker build --tag babylonlabs-io/ibcsim-bcd -f ibcsim-bcd/Dockerfile ${BABYLON_SDK_FULL_PATH}

ibcsim-bcd-rmi:
	docker rmi babylonlabs-io/ibcsim-bcd 2>/dev/null; true

start-bcd-bsn-integration: stop-bcd-bsn-integration
	ibcsim-bcd/pre-deployment.sh
	docker compose -f ibcsim-bcd/docker-compose.yml up -d

stop-bcd-bsn-integration:
	docker compose -f ibcsim-bcd/docker-compose.yml down
	rm -rf ibcsim-bcd/.testnets

.PHONY: ibcsim-bcd ibcsim-bcd-rmi start-bcd-bsn-integration stop-bcd-bsn-integration babylond babylond-rmi
