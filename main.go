package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

func usage() {
	name := os.Args[0]

	println("Convert CIDRs to glob-style patterns.")
	println("\nUsage:")
	printf("  %v [CIDR]...\n", name)
	printf("  %v -f file\n", name)
	printf("  command | %v\n", name)
	println("\nFlags:")
	flag.PrintDefaults()
}

func printf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func main() {
	var (
		file   string
		output string
	)

	flag.StringVar(&file, "f", "", "read from file instead of stdin")
	flag.StringVar(&output, "o", "", "write to file instead of stdout")
	flag.Usage = usage
	flag.Parse()

	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			println(err)
			return
		}

		defer func() {
			err := f.Close()
			if err != nil {
				println(err)
			}
		}()

		os.Stdout = f
	}

	if flag.NArg() > 0 {
		n := flag.NArg()
		for i := 0; i < n; i++ {
			parseCIDR(flag.Arg(i))
		}

		if file == "" {
			return
		}
	}

	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			println(err)
			return
		}
		defer f.Close()

		os.Stdin = f
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		if line != "" {
			parseCIDR(line)
		}
	}
}

func parseCIDR(cidr string) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		println("skip:", cidr)
		return
	}

	ip := ipnet.IP.To4()
	if ip == nil {
		println("skip:", cidr)
		return
	}

	ones, _ := ipnet.Mask.Size()
	if ones == 0 {
		println("skip:", cidr)
		return
	}

	if cidr != ipnet.String() {
		println("warning:", cidr, "changes to", ipnet)
	}

	switch {
	case ones <= 8:
		for _, s := range toGlob(int(ip[0]), int(ip[0])+1<<(8-ones)-1) {
			fmt.Printf("%v.[0-9]*.[0-9]*.[0-9]*\n", s)
		}
	case ones <= 16:
		for _, s := range toGlob(int(ip[1]), int(ip[1])+1<<(16-ones)-1) {
			fmt.Printf("%v.%v.[0-9]*.[0-9]*\n", int(ip[0]), s)
		}
	case ones <= 24:
		for _, s := range toGlob(int(ip[2]), int(ip[2])+1<<(24-ones)-1) {
			fmt.Printf("%v.%v.%v.[0-9]*\n", int(ip[0]), int(ip[1]), s)
		}
	case ones <= 32:
		for _, s := range toGlob(int(ip[3]), int(ip[3])+1<<(32-ones)-1) {
			fmt.Printf("%v.%v.%v.%v\n", int(ip[0]), int(ip[1]), int(ip[2]), s)
		}
	}
}

func toGlob(i, j int) (a []string) {
	if i == j {
		return []string{strconv.Itoa(i)}
	}

	if i > j {
		return
	}

	if i < 10 {
		if j < 10 {
			return []string{fmt.Sprintf("[%v-%v]", i, j)}
		}

		a = append(a, fmt.Sprintf("[%v-9]", i))
		i = 10
	}

	if x := i % 10; x > 0 {
		y := x + (j - i)
		if y < 10 {
			return append(a, fmt.Sprintf("%v[%v-%v]", i/10, x, y))
		}

		a = append(a, fmt.Sprintf("%v[%v-9]", i/10, x))
		i += 10 - x
	}

	var last string

	if y := j % 10; y < 9 {
		x := y - (j - i)
		if x >= 0 {
			return append(a, fmt.Sprintf("%v[%v-%v]", j/10, x, y))
		}

		last = fmt.Sprintf("%v[0-%v]", j/10, y)
		j -= y + 1
	}

	for _, s := range toGlob(i/10, j/10) {
		a = append(a, s+"[0-9]")
	}

	if last != "" {
		a = append(a, last)
	}

	return
}
