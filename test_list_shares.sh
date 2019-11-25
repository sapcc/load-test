DOMAIN_ID=${1?"domain id must be assigned as the first argument"}
PROJECT_ID=${2:?"project id must be assigned as the second argument"}

echo "get openstack token ..."
TOKEN=$(openstack --insecure token issue -c id -f value)
token_substring=$(echo $TOKEN | cut -c -30)

echo "domain  = $DOMAIN_ID"
echo "project = $PROJECT_ID"
echo "token   = ${token_substring}..."

header="X-Auth-Token: $TOKEN"
target="GET https://share-3.qa-de-1.cloud.sap/v2/${DOMAIN_ID}/shares?all_tenants=1&project_id=${PROJECT_ID}"

rate="${3-10}"
duration="${4-30s}"
output="${5-results.bin}"
echo "attack $target @$rate for $duration"

echo $target | vegeta attack -output=$output -header="$header" -keepalive=false -timeout=60s -duration=$duration -rate=$rate
