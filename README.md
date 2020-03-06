# load-test
http load tests on openstack endpoints based on  https://github.com/tsenart/vegeta

## manila
### pre-requisites
1. python-manilaclient, python-swiftclient
1. mkdir results
1. permission to CRUD shares and share snapshots in region qa-de-1:
  - domain_name = cc3test
  - project_name = zproject01_S4
1. permission to upload content to swift in region qa-de-1:
  - domain_name = monsoon3
  - project_name = storage_support
