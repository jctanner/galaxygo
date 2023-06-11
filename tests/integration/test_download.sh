#!/bin/bash

rm -rf /tmp/collections
ansible-galaxy collection install -vvv --no-cache -p /tmp/collections -s http://localhost:8080 ansible.posix
