# dwmsb

Output status bar metrics suitable for [dwm](https://dwm.suckless.org).

## Dependencies

### Files

`/proc/meminfo`, `/proc/loadavg`

### Commands

`amixer`, `df`, `hostname`, `iwgetid`

## Example

```v70% m38% 60/256 0.29/47 w:Home 10.0.0.72 2017-12-31 52 23:59```

* Volume 70%
* Memory usage 38%
* Disk usage 60GB of 254GB total
* Load (1-minute) 0.29, CPU temp 47 degrees
* Wireless connected to ESSID "Home"
* Local IP address is 10.0.0.72
* Date, week, time

## Installation

Set the X root window name regularly to the output of dwmsb. For example, add this to your `.xinitrc`:

```
while true
do
  xsetroot -name "$(/path/to/dwmsb)"
  sleep 1
done &
```
