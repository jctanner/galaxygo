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
        pulp_id::text,
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
        pulp_id::text,
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
        acv.content_ptr_id::text as pulp_id,
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
        acv.content_ptr_id::text as pulp_id,
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
        CONCAT('/api/v3/collections/', acv.namespace, '/', acv.name, '/') as collection_href,
        CONCAT('/api/v3/artifacts/', cca.relative_path) as download_url,
        acv.content_ptr_id::text as pulp_id,
        ac.pulp_id::text as collection_id,
        cc.pulp_created as created_at, 
        cc.pulp_last_updated as updated_at, 
        acv.collection_id::text,
        acv.namespace,
        acv.name,
        acv.version,
        acv.dependencies::text,
        acv.license,
        acv.repository,
        acv.is_highest,
        acv.requires_ansible,
        acv.description,
        acv.documentation,
        acv.homepage,
        acv.issues,
        acv.manifest,
        cca.artifact_id as artifact_id,
        cca.relative_path as filename,
        ca.sha256 as sha256,
        ca.size as size

    FROM ansible_collectionversion acv
    LEFT JOIN ansible_collection ac
        ON ac.namespace = acv.name AND ac.name = acv.name
    LEFT JOIN core_content cc
        ON cc.pulp_id = acv.content_ptr_id
    LEFT JOIN core_contentartifact cca
        ON cca.content_id = acv.content_ptr_id
    LEFT JOIN core_artifact ca
        ON ca.pulp_id = cca.artifact_id
    WHERE acv.namespace = '{{ namespace }}' AND acv.name = '{{ name }}' AND acv.version = '{{ version }}'
`

var ArtifactPathByFilename = `
    SELECT
        cca.relative_path as filename,
        ca.file as filepath
    FROM core_contentartifact cca
    LEFT JOIN core_artifact ca
        ON ca.pulp_id = cca.artifact_id
    WHERE cca.relative_path = '{{ filename }}'
    LIMIT 1
`


var NewArtifactObject = `
    INSERT INTO core_artifact
    (
        pulp_id,
        pulp_domain_id,
        pulp_created,
        file,
        size,
        sha256,
        timestamp_of_interest
    )
    VALUES
    (
        '{{ pulp_id }}',
        '{{ pulp_domain_id }}',
        '{{ pulp_created }}',
        '{{ file }}',
        '{{ size }}',
        '{{ sha256 }},
        '{{ timestamp_of_interest }}'
    )
`


var NewArtifactUploadTask = `
    INSERT INTO core_task
    (
        logging_cid,
        pulp_id,
        pulp_domain_id,
        pulp_created,
        state,
        name,
        args,
        kwargs,
    ) 
    VALUES
    (
        '{{ logging_cid }}',
        '{{ pulp_id }}',
        '{{ pulp_domain_id }}',
        '{{ pulp_created }}',
        'Waiting',
        'galaxy_ng.app.tasks.publishing.import_and_auto_approve',
        '{{ args }}',
        '{{ kwargs }}'
    )
`
