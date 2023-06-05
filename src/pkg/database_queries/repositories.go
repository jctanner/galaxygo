package database_queries

/*
pulp=# select * from core_repository where name='community';
-[ RECORD 1 ]--------+-------------------------------------
pulp_id              | 98a4bf6a-7009-4a62-a3f6-f97132cdf0d5
pulp_created         | 2023-06-04 15:05:08.427561+00
pulp_last_updated    | 2023-06-04 15:20:02.893428+00
name                 | community
description          | Community content repository
next_version         | 2
pulp_type            | ansible.ansible
remote_id            | 10e8e516-be56-4da8-9463-8d39dbc450c7
retain_repo_versions | 1
user_hidden          | f
pulp_labels          |
pulp_domain_id       | 0cc6d6ba-e465-4584-b404-824eba1d0a3a
*/

var GetRepositoryIdByName = `
    SELECT
		pulp_id::text,
		name
    FROM core_repository
	WHERE name='{{ repository_name }}'
`
