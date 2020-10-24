clean:
	rm -f nft-exchangecli
	rm -f nft-exchanged
install:
	go install -mod=readonly ./cmd/nft-exchangecli
	go install -mod=readonly ./cmd/nft-exchanged

