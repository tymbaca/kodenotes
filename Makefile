run:
	export TARGET_STAGE=run &&\
	docker compose up --build
test:
	export TARGET_STAGE=run &&\
	docker compose up --build
