package database_queries

var CheckUsernameAndPassword = `
    SELECT
		count(id)
    FROM galaxy_user
    WHERE username='{{ username }}' AND password='{{ password }}'
`
