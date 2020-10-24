clean:
	rm -f nftchaincli
	rm -f nftchaind
install:
	go install ./cmd/nftchaincli
	go install ./cmd/nftchaind
