#!/usr/bin/env bash
DOMAIN_NAME=${1-monsoon3}
# storage_support project
PROJECT_ID=${2-05f9781218b7401d9955f9b8a05a5aea}
CONTAINER=${3-manila-load-test}

TOKEN=$(openstack --os-user-domain-name=$DOMAIN_NAME --os-project-id=$PROJECT_ID --insecure token issue -c id -f value)
# TODO: use attack date, not current date
# ideally the script would get idempotent
DATE=$(date +%Y-%m-%d_%H-%M)


echo "uploading to https://objectstore-3.qa-de-1.cloud.sap/v1/AUTH_$PROJECT_ID/$CONTAINER/"

for file in results/*.{html,txt}
do
  swift --os-auth-token $TOKEN upload $CONTAINER $file --object-name ${DATE}_$file
done
