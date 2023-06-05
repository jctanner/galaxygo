package database_queries

type NamespacePostData struct {
	Name string `json:"name"`
}

var ListNamespaces = `
    SELECT
		id,
		name
    FROM galaxy_namespace
`

var CheckIfNamespaceExists = `
    SELECT
        count(id)
    FROM galaxy_namespace
	WHERE name='{{ name }}'
`

var GetNamespaceByName = `
    SELECT
        id,
		name,
		company,
		email,
		_avatar_url as avatar_url,
		description
    FROM galaxy_namespace
	WHERE name='{{ namespace_name }}'
`

var CreateNamespace = `
	INSERT INTO galaxy_namespace
		(
			name,
			company,
			email,
			_avatar_url,
			description,
			resources
		)
	VALUES
		(
			'{{ namespace_name }}',
			'',
			'',
			'',
			'',
			''
		)
`
