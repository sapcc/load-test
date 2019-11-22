
share:
	rm -f output.txt
	sleep 5 && ./list_shares.sh 10 5s output.txt
	sleep 5 && ./list_shares.sh 20 5s output.txt
	sleep 5 && ./list_shares.sh 30 5s output.txt
