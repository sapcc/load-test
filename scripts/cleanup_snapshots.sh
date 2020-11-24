DOMAIN_ID=${1?"domain id must be assigned as the first argument"}
PROJECT_ID=${2:?"project id must be assigned as the second argument"}

echo "get openstack token for snapshot cleanup ..."
TOKEN=$(openstack --os-user-domain-id=$DOMAIN_ID --os-project-domain-id=$DOMAIN_ID --os-project-id=$PROJECT_ID --insecure token issue -c id -f value)
token_substring=$(echo $TOKEN | cut -c -30)
URL="https://share-3.qa-de-1.cloud.sap/v2/${PROJECT_ID}"

echo "domain  = $DOMAIN_ID"
echo "project = $PROJECT_ID"
echo "token   = ${token_substring}..."
echo "url     = $URL"

for snap in $(manila --os-user-domain-id=$DOMAIN_ID --os-project-id=$PROJECT_ID --bypass-url $URL --os-token $TOKEN snapshot-list | grep -E 'available' | awk '{ print $2 }')
do
  manila --os-user-domain-id=$DOMAIN_ID --os-project-id=$PROJECT_ID --bypass-url $URL --os-token $TOKEN snapshot-force-delete $snap
done
