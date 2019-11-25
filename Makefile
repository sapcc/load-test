domain=9e2ae21cd643430f8efe9005a758c4e7
project=b25c933351f54df5a9a122b72ba523c8

sleeptime=60
attacktime=5m

all:

share:
	@rm -f results/results-10-$(attacktime).bin
	@rm -f results/results-30-$(attacktime).bin
	@rm -f results/results-50-$(attacktime).bin
	@echo "sleep $(sleeptime) sec" && sleep $(sleeptime)
	./test_list_shares.sh $(domain) $(project) 10 $(attacktime) results/results-10-$(attacktime).bin
	@echo "sleep $(sleeptime) sec" && sleep $(sleeptime)
	./test_list_shares.sh $(domain) $(project) 30 $(attacktime) results/results-30-$(attacktime).bin
	@echo "sleep $(sleeptime) sec" && sleep $(sleeptime)
	./test_list_shares.sh $(domain) $(project) 50 $(attacktime) results/results-50-$(attacktime).bin

snapshot: manila-load-test/manila-load-test
	./test_create_snapshots.sh $(domain) $(project) 1 15s

manila-load-test/manila-load-test: manila-load-test/*.go
	cd manila-load-test && make

.Phony: shares
shares: 
	manila list --all-tenant --project-id $(project) | cut -d\| -f 2 | xargs -n 1 echo | tail -n +4 | head --lines=-1 > $@.txt
