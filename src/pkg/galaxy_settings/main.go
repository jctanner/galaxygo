package galaxy_settings

import (
	"os"
	"strconv"
)

type GalaxySettings struct {
	Analytics bool

	Api_hostname string
	Api_prefix   string

	Content_hostname    string
	Content_path_prefix string
	Content_origin      string

	Default_distribution_base_path string

	Deployment_mode string

	Token_auth_disabled       bool
	Token_server              string
	Token_signature_algorithm string

	Public_key_path  string
	Private_key_path string

	Require_content_approval   string
	Auto_sign_collections      bool
	Collection_signing_service string
	Rh_entitlement_required    string

	Enable_unauthenticated_collection_access   bool
	Enable_unauthenticated_collection_download bool
	Enable_legacy_roles                        bool
	Enable_execution_environments              bool

	Display_repositories bool
	Ai_deny_index        bool

	Gunicorn_workers int
	Debug            bool

	Use_redis      bool
	Redis_address  string
	Redis_password string
	Redis_db       int

	Use_s3              bool
	Aws_access_key      string
	Aws_secret_key      string
	Aws_s3_region       string
	Aws_s3_endpoint_url string
	Aws_s3_bucket_name  string

	Default_file_storage string

	Database_host     string
	Database_name     string
	Database_username string
	Database_password string

	Social_auth_login_redirect_url string
}

func get_env_list(env_vars []string, fallback string) string {
	for _, env_var := range env_vars {
		value := os.Getenv(env_var)
		if value != "" {
			return value
		}
	}
	return fallback
}

func StringToBool(str string) bool {
	if str == "" {
		return false
	}

	b, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}

	return b
}

func NewGalaxySettings() GalaxySettings {

	debug := StringToBool(get_env_list([]string{"DEBUG"}, "true"))

	api_prefix := get_env_list([]string{"PULP_GALAXY_API_PATH_PREFIX"}, "")

	default_distribution_base_path := get_env_list(
		[]string{"GALAXY_API_DEFAULT_DISTRIBUTION_BASE_PATH"},
		"community",
	)

	database_host := get_env_list(
		[]string{"DATABASE_HOST", "PULP_DATABASES__default__HOST"},
		"",
	)
	database_name := get_env_list(
		[]string{"DATABASE_NAME", "POSTGRES_DB", "PULP_DATABASES__default__NAME"},
		"",
	)
	database_user := get_env_list(
		[]string{"DATABASE_USER", "POSTGRES_USER", "PULP_DATABASES__default__USER"},
		"",
	)
	database_pass := get_env_list(
		[]string{"DATABASE_PASSWORD", "POSTGRES_PASSWORD", "PULP_DATABASES__default__PASSWORD"},
		"",
	)

	use_redis := StringToBool(get_env_list([]string{"USE_REDIS"}, "false"))
	redis_address := get_env_list([]string{"REDIS_ADDRESS"}, "redis:6379")
	redis_password := get_env_list([]string{"REDIS_PASSWORD"}, "")
	redis_db, _ := strconv.Atoi(get_env_list([]string{"REDIS_DATABASE"}, "0"))

	aws_access_key := get_env_list([]string{"PULP_AWS_ACCESS_KEY_ID"}, "")
	aws_secret_key := get_env_list([]string{"PULP_AWS_SECRET_ACCESS_KEY"}, "")
	aws_s3_region := get_env_list([]string{"PULP_AWS_S3_REGION_NAME"}, "")
	aws_s3_endpoint := get_env_list([]string{"PULP_AWS_S3_ENDPOINT_URL"}, "")
	aws_s3_bucket := get_env_list([]string{"PULP_AWS_STORAGE_BUCKET_NAME"}, "")
	default_file_storage := get_env_list([]string{"PULP_DEFAULT_FILE_STORAGE"}, "")

	use_s3 := false
	if default_file_storage == "storages.backends.s3boto3.S3Boto3Storage" {
		use_s3 = true
	}

	return GalaxySettings{
		Api_prefix:                     api_prefix,
		Debug:                          debug,
		Use_redis:                      use_redis,
		Redis_address:                  redis_address,
		Redis_password:                 redis_password,
		Redis_db:                       redis_db,
		Aws_access_key:                 aws_access_key,
		Aws_secret_key:                 aws_secret_key,
		Aws_s3_region:                  aws_s3_region,
		Aws_s3_endpoint_url:            aws_s3_endpoint,
		Aws_s3_bucket_name:             aws_s3_bucket,
		Use_s3:                         use_s3,
		Database_host:                  database_host,
		Database_name:                  database_name,
		Database_username:              database_user,
		Database_password:              database_pass,
		Default_distribution_base_path: default_distribution_base_path,
	}
}
