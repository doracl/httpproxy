/* IP
 */

package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}

	name := os.Args[1]

	// addr := net.ParseIP(name)
	// if addr == nil {
	// 	fmt.Println("Invalid address")
	// } else {
	// 	fmt.Println("The address is ", addr.String())
	// }

	// mask := addr.DefaultMask()
	// network := addr.Mask(mask)
	// ones, bits := mask.Size()
	// fmt.Println("Address is ", addr.String(), " Defautl mask length is ", bits,
	// 	"Leading ones count is ", ones, "Mask is (hex) ", mask.String(), "Network is ", network.String())

	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		fmt.Println("Resolution error ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Resolved address is ", addr.String())

	addrs, err := net.LookupHost(name)

	if err != nil {
		fmt.Print("Err: ", err.Error())
		os.Exit(2)
	}

	for _, s := range addrs {
		fmt.Println(s)
	}
	os.Exit(0)
}
