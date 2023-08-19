run: 
	set -o allexport && . ./.env && set +o allexport &&\
	export POSTGRES_HOST=localhost &&\
	go run main.go

test: 
	set -o allexport && . ./.env && set +o allexport &&\
	export POSTGRES_HOST=localhost &&\
	go test -v ./... -count=1

pg-up:
	set -o allexport && . ./.env && set +o allexport &&\
	docker run --name kodenotes-postgres -e POSTGRES_PASSWORD=$${POSTGRES_PASSWORD} -p 5432:5432 -d postgres

pg-down:
	docker stop kodenotes-postgres &&\
	docker rm kodenotes-postgres 

pg-restart:
	set -o allexport && . ./.env && set +o allexport &&\
	docker stop kodenotes-postgres &&\
	docker rm kodenotes-postgres &&\
	docker run --name kodenotes-postgres -e POSTGRES_PASSWORD=$${POSTGRES_PASSWORD} -p 5432:5432 -d postgres

pg-shell:
	docker exec -it kodenotes-postgres psql -U postgres -d postgres
