run: 
	export POSTGRES_HOST=localhost &&\
	export POSTGRES_PASSWORD=mypassword &&\
	go run main.go

test:
	export POSTGRES_HOST=localhost &&\
	export POSTGRES_PASSWORD=mypassword &&\
	go test ./...

compose-run:
	export TARGET_STAGE=run &&\
	docker compose up --build

compose-test:
	export TARGET_STAGE=test &&\
	docker compose up --build

pg-up:
	docker run --name kodenotes-postgres -e POSTGRES_PASSWORD=mypassword -p 5432:5432 -d postgres

pg-down:
	docker stop kodenotes-postgres 
	docker rm kodenotes-postgres 

pg-restart:
	docker stop kodenotes-postgres 
	docker rm kodenotes-postgres 
	docker run --name kodenotes-postgres -e POSTGRES_PASSWORD=mypassword -p 5432:5432 -d postgres

pg-shell:
	docker exec -it kodenotes-postgres psql -U postgres -d postgres
