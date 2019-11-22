DOMAIN_ID=9e2ae21cd643430f8efe9005a758c4e7
PROJECT_ID=b25c933351f54df5a9a122b72ba523c8

echo "get openstack token ..."
TOKEN=$(openstack --insecure token issue -c id -f value)
token_substring=$(echo $TOKEN | cut -c -30)

echo "domain  = $DOMAIN_ID"
echo "project = $PROJECT_ID"
echo "token   = ${token_substring}..."

header="X-Auth-Token: $TOKEN"
target="GET https://share-3.qa-de-1.cloud.sap/v2/${DOMAIN_ID}/shares/detail?all_tenants=1&project_id=${PROJECT_ID}"

rate="${1-10}"
duration="${2-30s}"
output="${3-output.txt}"
echo "attack $target @$rate per seconds for $duration"

echo $target | \
  vegeta attack -header="$header" -duration=$duration -rate=$rate | \
  tee results.bin | \
  vegeta report --type json \
  | jq >> $output

#echo $target | \
#  vegeta attack -header="$header" -duration=$duration -rate=$rate --output output.txt | \
#  tee results.bin | \
#  vegeta report

