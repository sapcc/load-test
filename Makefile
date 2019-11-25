domain=9e2ae21cd643430f8efe9005a758c4e7
project=b25c933351f54df5a9a122b72ba523c8

sleeptime=60
attacktime=5m

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

shares.txt:
	manila list --all-tenant --project-id $(project) | cut -d\| -f 2 > $@
