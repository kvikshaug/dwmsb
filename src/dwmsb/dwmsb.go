package main

import (
	// "bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func battery() (s string) {
	out, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/uevent")
	check(err)

	lines := strings.Split(string(out), "\n")
	char := ""

	if strings.Contains(lines[1], "Full") {
		char = "f"
	} else if strings.Contains(lines[1], "Discharging") {
		char = "d"
	} else if strings.Contains(lines[1], "Charging") {
		char = "c"
	} else {
		char = "?"
	}

	full, err := strconv.ParseFloat(strings.Split(lines[9], "=")[1], 64)
	check(err)

	now, err := strconv.ParseFloat(strings.Split(lines[10], "=")[1], 64)
	check(err)

	pct := int64(((now / full) * 100) + 0.5)
	return fmt.Sprintf("%s%d", char, pct)
}

func audio() (s string) {
	out, err := exec.Command("amixer", "get", "Master").Output()
	check(err)

	lines := strings.Split(string(out), "\n")

	volume := "?"
	for _, line := range lines {
		if !strings.Contains(line, "[") {
			continue
		}

		if strings.Contains(line, "[off]") {
			volume = "M"
			break
		} else {
			r := strings.Trim(strings.Fields(line)[3], "[]")
			volume = fmt.Sprintf("v%s", r)
			break
		}
	}

	return volume
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
	if (pct > 80) {
		return fmt.Sprintf("m%d%% ** WARNING: HIGH MEMORY USAGE **", pct)
	} else {
		return fmt.Sprintf("m%d%%", pct)
	}
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

	matches, err := filepath.Glob("/sys/class/thermal/thermal_zone*")
	check(err)
	temp := "?"
	for index := range matches {
		out, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/thermal/thermal_zone%d/type", index))
		check(err)
		if strings.Trim(string(out), "\n") == "x86_pkg_temp" {
			out, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/thermal/thermal_zone%d/temp", index))
			check(err)
			tempMillis, err := strconv.ParseFloat(strings.Trim(string(out), "\n"), 64)
			check(err)
			temp = fmt.Sprintf("%.0f", tempMillis / 1000)
		}
	}

	return fmt.Sprintf("%s/%s", load, temp)
}

func network() (s string) {
	out, err := exec.Command("iwgetid").Output()
	var iwname string
	if err != nil {
		iwname = "↓"
	} else {
		s = string(out)
		a := strings.Index(s, "\"")
		b := strings.LastIndex(s, "\"")
		iwname = s[a+1 : b]
	}

	files, err := ioutil.ReadDir("/sys/class/net")
	check(err)

	ethstatus := "?"
	for _, folder := range files {
		if strings.HasPrefix(folder.Name(), "enp") {
			out, err = ioutil.ReadFile(fmt.Sprintf("/sys/class/net/%s/operstate", folder.Name()))
			check(err)

			if (string(out) == "up\n") {
				ethstatus = "↑"
			} else {
				ethstatus = "↓"
			}
		}
	}

	return fmt.Sprintf("w:%s e:%s", iwname, ethstatus)
}

func date() (s string) {
	now := time.Now()
	_, week := now.ISOWeek()
	return fmt.Sprintf("%s %d %s", now.Format("2006-01-02"), week, now.Format("15:04"))
}

func main() {
	l := []string{
		battery(),
		audio(),
		memory(),
		disk_usage(),
		cpu(),
		network(),
		date(),
	}
	fmt.Printf("%s", strings.Join(l, " "))
}
