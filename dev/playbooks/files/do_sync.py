#!/usr/bin/env python3

import requests
import yaml
import sys
import time

from logzero import logger

# http://localhost:8002/api/pulp/api/v3/remotes/ansible/collection/d7daddfb-b989-4011-b3bc-62eebe660a36/
'''
{
    "name":"community",
    "url":"https://beta-galaxy.ansible.com/api/",
    "tls_validation":true,
    "download_concurrency":10,
    "rate_limit":8,
    "requirements_file":"# Sample requirements.yaml\n\ncollections:\n  - name: ansible.posix\n  - name: community.network",
    "signed_only":false,
    "hidden_fields":[
        {"name":"client_key","is_set":false},
        {"name":"proxy_username","is_set":false},
        {"name":"proxy_password","is_set":false},
        {"name":"username","is_set":false},
        {"name":"password","is_set":false},
        {"name":"token","is_set":false}
    ],
    "pulp_href":"/api/pulp/api/v3/remotes/ansible/collection/d7daddfb-b989-4011-b3bc-62eebe660a36/",
    "pulp_created":"2023-05-26T22:55:38.254022Z",
    "pulp_labels":{},
    "pulp_last_updated":"2023-05-26T22:55:45.545674Z",
    "policy":"immediate",
    "sync_dependencies":true
}
'''


def main():

    #baseurl = 'http://localhost:5001'
    baseurl = sys.argv[1]

    remotes_rr = requests.get(
        baseurl + '/api/pulp/api/v3/repositories/ansible/ansible/',
        auth=('admin', 'admin')
    )
    repositories_map = dict((x['name'], x) for x in remotes_rr.json()['results'])
    community_repo_url = baseurl + repositories_map['community']['pulp_href']

    remotes_rr = requests.get(
        baseurl + '/api/pulp/api/v3/remotes/ansible/collection/',
        auth=('admin', 'admin')
    )
    remotes_map = dict((x['name'], x) for x in remotes_rr.json()['results'])
    community_remote_url = baseurl + remotes_map['community']['pulp_href']

    spec = {
        'collections': [
            'ansible.posix',
            'community.network'
        ]
    }
    spec_yaml = yaml.dump(spec)
    payload = {
        "name": "community",
        "url": "https://beta-galaxy.ansible.com/api/",
        'requirements_file': spec_yaml
    }
    rr = requests.put(
        community_remote_url,
        json=payload,
        auth=('admin', 'admin')
    )
    task_url = baseurl + rr.json()['task']
    while True:
        rrt = requests.get(task_url, auth=('admin', 'admin'))
        state = rrt.json()['state']
        logger.info(state)
        if state in ['completed', 'failed']:
            break
        time.sleep(2)

    rr = requests.post(
        community_repo_url.rstrip('/') + '/sync/',
        auth=('admin', 'admin')
    )
    task_url = baseurl + rr.json()['task']
    while True:
        rrt = requests.get(task_url, auth=('admin', 'admin'))
        state = rrt.json()['state']
        logger.info(state)
        if state in ['completed', 'failed']:
            break
        time.sleep(5)


if __name__ == "__main__":
    main()
