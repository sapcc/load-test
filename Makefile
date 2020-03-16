domain_name=cc3test
project_name=zproject01_S4
domain_id=603432c20d184a5e9699900914ad7765
project_id=b25c933351f54df5a9a122b72ba523c8

sleeptime=60
attacktime=5m

all:

share:
	@rm -f results/share-list-10-$(attacktime).bin
	@rm -f results/share-list-30-$(attacktime).bin
	@rm -f results/share-list-50-$(attacktime).bin
	./test_list_shares.sh $(domain_id) $(project_id) 10 $(attacktime) results/share-list-10-$(attacktime).bin
	@echo "sleep $(sleeptime) sec" && sleep $(sleeptime)
	./test_list_shares.sh $(domain_id) $(project_id) 30 $(attacktime) results/share-list-30-$(attacktime).bin
	@echo "sleep $(sleeptime) sec" && sleep $(sleeptime)
	./test_list_shares.sh $(domain_id) $(project_id) 50 $(attacktime) results/share-list-50-$(attacktime).bin

snapshot: manila-load-test/manila-load-test
	./test_create_snapshots.sh $(domain_id) $(project_id) 3 30s
	@echo "sleep $(sleeptime) sec" && sleep $(sleeptime)
	./cleanup_snapshots.sh $(domain_id) $(project_id)

manila-load-test/manila-load-test: manila-load-test/*.go
	cd manila-load-test && make

public:
	./generate_output.sh
	./upload_results_to_swift.sh

.Phony: shares
shares:
	manila --os-user-domain-name=$(domain_name) --os-project-domain-name=$(domain_name) --os-project-name=$(project_name) list | grep -E 'NFS' | awk '{ print $$2 }' > $@.txt
