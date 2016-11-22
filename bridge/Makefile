# Old-skool build tools.
#
OUT_DIR = _output
BIN_NAME = ipam
all build:
	mkdir -p ${OUT_DIR}/bin
	go build -o ${OUT_DIR}/bin/ipam main/main.go
.PHONY: all build

clean:
	rm -rf $(OUT_DIR)
.PHONY: clean
