#!/bin/bash

export DATABASE_HOST=$(docker inspect galaxy_ng-postgres-1 | jq '.[0].NetworkSettings.Networks.galaxy_ng_default.IPAddress' | tr -d '"')
export DATABASE_NAME=$(docker inspect galaxy_ng-postgres-1 | jq '.[0].Config.Env' | fgrep POSTGRES_DB | cut -d\= -f2 | tr -d '"' | tr -d ',')
export DATABASE_USER=$(docker inspect galaxy_ng-postgres-1 | jq '.[0].Config.Env' | fgrep POSTGRES_USER | cut -d\= -f2 | tr -d '"' | tr -d ',')
export DATABASE_PASSWORD=$(docker inspect galaxy_ng-postgres-1 | jq '.[0].Config.Env' | fgrep POSTGRES_PASSWORD | cut -d\= -f2 | tr -d '"' | tr -d ',')

source extra_settings.sh

go run api/api.go
