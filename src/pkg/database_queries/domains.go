package database_queries

var GetDefaultDomainID = `
	SELECT
		pulp_id::text
	FROM
		core_domain
	WHERE
		name = 'default'
`
