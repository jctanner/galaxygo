package database_queries

var CheckUsernameAndPassword = `
    SELECT
		count(id)
    FROM galaxy_user
    WHERE username='{{ username }}' AND password='{{ password }}'
`

var GetUsernameAndPassword = `
    SELECT
        id,
        username,
        password
    FROM galaxy_user
    WHERE username='{{ username }}'
`
