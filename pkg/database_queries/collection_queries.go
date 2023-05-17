package database_queries

/*************************************************************************

    Collections

*************************************************************************/

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
        name,
        CONCAT('/api/v3/collections/', namespace, '/', name, '/') as href
    FROM ansible_collection
`

var CollectionSummary = `
    SELECT
        CONCAT('/api/v3/collections/', namespace, '/', name, '/') as href,
        CONCAT('/api/v3/collections/', namespace, '/', name, '/versions/') as versions_url,
        pulp_id,
        pulp_created as created_at,
        pulp_last_updated as updated_at,
        namespace,
        name
    FROM ansible_collection
    WHERE namespace='{{ namespace }}' AND name='{{ name }}'
`

var CollectionVersionsSummary = `
    SELECT
        CONCAT('/api/v3/collections/', acv.namespace, '/', acv.name, '/versions/', acv.version, '/') as href,
        acv.content_ptr_id as pulp_id,
        cc.pulp_created as created_at,
        cc.pulp_last_updated as updated_at,
        acv.version,
        acv.requires_ansible
    FROM ansible_collectionversion acv
    LEFT JOIN core_content cc
        ON cc.pulp_id = acv.content_ptr_id
    WHERE acv.namespace='{{ namespace }}' AND acv.name='{{ name }}'
`

var CollectionVersionsSummaryCount = `
    SELECT
        count(content_ptr_id)
    FROM ansible_collectionversion
    WHERE namespace='{{ namespace }}' AND name='{{ name }}'
`

/*************************************************************************

    Collection Versions

*************************************************************************/

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

var CollectionVersionDetail = `
    SELECT
        CONCAT('/api/v3/collections/', acv.namespace, '/', acv.name, '/versions/', acv.version, '/') as href,
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
        acv.requires_ansible,
        acv.description,
        acv.documentation,
        acv.homepage,
        acv.issues,
        acv.manifest
    FROM ansible_collectionversion acv
    LEFT JOIN core_content cc
        ON cc.pulp_id = acv.content_ptr_id
    WHERE acv.namespace = '{{ namespace }}' AND name = '{{ name }}' AND version = '{{ version }}'
`
