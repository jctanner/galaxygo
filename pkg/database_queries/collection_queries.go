package database_queries

var CountCollections = `
    SELECT
        count(pulp_id)
    FROM ansible_collection
`

var ListCollections = `
    SELECT
        pulp_id,
        pulp_created as created_at,
        pulp_last_updated as updated_at,
        namespace,
        name
    FROM ansible_collection
`

var CountCollectionVersions = `
    SELECT
        count(content_ptr_id)
    FROM ansible_collectionversion
`

var ListCollectionVersions = `
    SELECT
        acv.content_ptr_id as pulp_id,
        cc.pulp_created as created_at, 
        cc.pulp_last_updated as updated_at, 
        acv.collection_id,
        acv.namespace,
        acv.name,
        acv.version,
        acv.dependencies,
        acv.license,
        acv.repository,
        acv.is_highest,
        acv.requires_ansible
    FROM ansible_collectionversion acv
    LEFT JOIN core_content cc
        ON cc.pulp_id = acv.content_ptr_id
`