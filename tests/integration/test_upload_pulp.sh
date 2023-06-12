#!/bin/bash

WORKDIR=$(mktemp -d)

cd $WORKDIR

echo "[galaxy]" > ansible.cfg
echo "server_list=golang" >> ansible.cfg
echo "[galaxy_server.golang]" >> ansible.cfg
#echo "url=http://localhost:8080" >> ansible.cfg
echo "url=http://localhost:5001" >> ansible.cfg
echo "username=admin" >> ansible.cfg
echo "password=admin" >> ansible.cfg

cat ansible.cfg

ansible-galaxy collection init foo.bar

#find .

cd foo/bar
ansible-galaxy collection build .

cd $WORKDIR

find .

curl -u admin:admin -X POST -H 'Content-Type: application/json' -d '{"name": "foo", "groups": []}' http://localhost:8080/api/_ui/v1/namespaces/

ansible-galaxy collection publish -vvvv foo/bar/foo-bar-*.tar.gz
