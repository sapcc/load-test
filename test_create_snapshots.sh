#!/usr/bin/env bash
DOMAIN_ID=${1?"domain id must be assigned as the first argument"}
PROJECT_ID=${2:?"project id must be assigned as the second argument"}
rate="${3-5}"
duration="${4-1s}"

echo "get openstack token for snapshot create ..."
TOKEN=$(openstack --os-user-domain-id=$DOMAIN_ID --os-project-domain-id=$DOMAIN_ID --os-project-id=$PROJECT_ID --insecure token issue -c id -f value)
token_substring=$(echo $TOKEN | cut -c -30)
url="https://share-3.qa-de-1.cloud.sap/v2/${PROJECT_ID}"

echo "domain  = $DOMAIN_ID"
echo "project = $PROJECT_ID"
echo "token   = ${token_substring}..."
echo "url     = $url"

./manila-load-test/manila-load-test -url $url -token $TOKEN -shares ./shares.txt -rate $rate -duration $duration > ./results/snap-create-$rate-$duration.bin
