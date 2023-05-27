#!/bin/bash

export DATABASE_HOST=$(docker inspect galaxy_ng-postgres-1 | jq '.[0].NetworkSettings.Networks.galaxy_ng_default.IPAddress' | tr -d '"')
export DATABASE_NAME=$(docker inspect galaxy_ng-postgres-1 | jq '.[0].Config.Env' | fgrep POSTGRES_DB | cut -d\= -f2 | tr -d '"' | tr -d ',')
export DATABASE_USER=$(docker inspect galaxy_ng-postgres-1 | jq '.[0].Config.Env' | fgrep POSTGRES_USER | cut -d\= -f2 | tr -d '"' | tr -d ',')
export DATABASE_PASSWORD=$(docker inspect galaxy_ng-postgres-1 | jq '.[0].Config.Env' | fgrep POSTGRES_PASSWORD | cut -d\= -f2 | tr -d '"' | tr -d ',')


export PULP_CONTENT_PATH_PREFIX=/api/v3/artifacts/collections/

export PULP_GALAXY_API_PATH_PREFIX=/api/
export PULP_GALAXY_DEPLOYMENT_MODE=standalone
export PULP_GALAXY_REQUIRE_CONTENT_APPROVAL=false
export PULP_GALAXY_AUTO_SIGN_COLLECTIONS=false
export PULP_GALAXY_COLLECTION_SIGNING_SERVICE=ansible-default
export PULP_RH_ENTITLEMENT_REQUIRED=insights
export PULP_ANSIBLE_API_HOSTNAME=http://localhost:5001
export PULP_ANSIBLE_CONTENT_HOSTNAME=http://localhost:24816/api/v3/artifacts/collections
export PULP_TOKEN_AUTH_DISABLED=true
export PULP_CONTENT_ORIGIN="http://localhost:24816"
export PULP_GALAXY_FEATURE_FLAGS__display_repositories=false
export PULP_TOKEN_AUTH_DISABLED=false
export PULP_TOKEN_SERVER=http://localhost:5001/token/
export PULP_TOKEN_SIGNATURE_ALGORITHM=ES256
export PULP_PUBLIC_KEY_PATH=/src/galaxy_ng/dev/common/container_auth_public_key.pem
export PULP_PRIVATE_KEY_PATH=/src/galaxy_ng/dev/common/container_auth_private_key.pem

export PULP_ANALYTICS=false
export PULP_GALAXY_ENABLE_UNAUTHENTICATED_COLLECTION_ACCESS=true
export PULP_GALAXY_ENABLE_UNAUTHENTICATED_COLLECTION_DOWNLOAD=true

#export SOCIAL_AUTH_GITHUB_KEY="${SOCIAL_AUTH_GITHUB_KEY:null}"
#export SOCIAL_AUTH_GITHUB_SECRET="${SOCIAL_AUTH_GITHUB_SECRET:null}"

export PULP_GALAXY_ENABLE_LEGACY_ROLES=true
export PULP_GALAXY_FEATURE_FLAGS__execution_environments=false
export PULP_SOCIAL_AUTH_LOGIN_REDIRECT_URL="/"
export PULP_GALAXY_FEATURE_FLAGS__ai_deny_index=true

####################################
# MINIO
####################################
export PULP_DEFAULT_FILE_STORAGE="storages.backends.s3boto3.S3Boto3Storage"
export PULP_MEDIA_ROOT=""
#PULP_AWS_ACCESS_KEY_ID="AKIAIT2Z5TDYPX3ARJBA"
#PULP_AWS_SECRET_ACCESS_KEY="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"
export PULP_AWS_ACCESS_KEY_ID="minioadmin"
export PULP_AWS_SECRET_ACCESS_KEY="minioadmin"
export PULP_AWS_S3_REGION_NAME="eu-central-1"
export PULP_AWS_S3_ADDRESSING_STYLE="path"
export PULP_S3_USE_SIGV4=true
export PULP_AWS_S3_SIGNATURE_VERSION="s3v4"
export PULP_AWS_STORAGE_BUCKET_NAME="pulp3"
export PULP_AWS_S3_ENDPOINT_URL="http://minio:9000"
export PULP_AWS_DEFAULT_ACL="@none None"

cd src
go run api/api.go
