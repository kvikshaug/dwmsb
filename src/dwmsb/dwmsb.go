package main

import (
	// "bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func audio() (s string) {
	out, err := exec.Command("amixer", "get", "Master").Output()
	check(err)

	lines := strings.Split(string(out), "\n")
	line := lines[4]

	if strings.Contains(lines[4], "[off]") {
		return "M"
	} else {
		r := strings.Trim(strings.Fields(line)[3], "[]")
		return fmt.Sprintf("v%s", r)
	}
}

func memory() (s string) {
	out, err := ioutil.ReadFile("/proc/meminfo")
	check(err)

	lines := strings.Split(string(out), "\n")

	total, err := strconv.ParseFloat(strings.Fields(lines[0])[1], 64)
	check(err)

	available, err := strconv.ParseFloat(strings.Fields(lines[2])[1], 64)
	check(err)

	pct := int64(100 - ((available / total) * 100) + 0.5)
	return fmt.Sprintf("m%d%%", pct)
}

func disk_usage() (s string) {
	out, err := exec.Command("df", "-h", ".").Output()
	check(err)

	lines := strings.Split(string(out), "\n")
	line := lines[1]

	values := strings.Fields(line)
	return fmt.Sprintf("%s/%s", strings.Trim(values[2], "G"), strings.Trim(values[1], "G"))
}

func cpu() (s string) {
	out, err := ioutil.ReadFile("/proc/loadavg")
	check(err)

	load := strings.Split(string(out), "\n")[0][:4]

	out, err = exec.Command("sensors").Output()
	check(err)

	line := strings.Split(string(out), "\n")[2]
	temp_full := strings.Fields(line)[1]
	temp := temp_full[1:len(temp_full) - 5]

	return fmt.Sprintf("%s/%s", load, temp)
}

func network() (s string) {
	// out, err := exec.Command("ip", "-o", "addr", "show", "up", "primary", "scope", "global").Output()
	// check(err)

	// lines := strings.Split(string(out), "\n")
	// var output bytes.Buffer
	// for _, line := range lines {
	// 	if len(line) == 0 {
	// 		continue
	// 	}
	// 	columns := strings.Fields(line)
	// 	address := columns[3]
	// 	address = address[:len(address)-3]
	// 	output.WriteString(fmt.Sprintf("%s ", address))
	// }

	// addresses := output.String()
	// addresses = addresses[:len(addresses)-1]

	out, err := exec.Command("hostname", "-i").Output()
	check(err)
	addresses := strings.Trim(string(out), "\n ")

	out, err = exec.Command("iwgetid").Output()
	var iwname string
	if err != nil {
		iwname = "↓"
	} else {
		s = string(out)
		a := strings.Index(s, "\"")
		b := strings.LastIndex(s, "\"")
		iwname = s[a+1:b]
	}
	return fmt.Sprintf("w:%s %s", iwname, addresses)
}

func date() (s string) {
	now := time.Now()
	_, week := now.ISOWeek()
	return fmt.Sprintf("%s %d %s", now.Format("2006-01-02"), week, now.Format("15:04"))
}

func main() {
	l := []string {
		audio(),
		memory(),
		disk_usage(),
		cpu(),
		network(),
		date(),
	}
	fmt.Printf("%s", strings.Join(l, " "))
}
