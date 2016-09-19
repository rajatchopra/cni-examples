
package main

import (
	"github.com/rajatchopra/cni-examples/bridge/ipam"
	"fmt"
	"bufio"
	"encoding/json"
	"os"
)

func main() {
	ipam, err := ipam.NewIPAM("")
	if err != nil {
		fmt.Printf("ERROR[1]: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		reader.ReadString('\n')
		ip, err := ipam.GetNextIP()
		if err != nil {
			fmt.Printf("ERROR[2]: %v", err)
			return
		}
		fmt.Println(ip.String())
	}
	return
}
