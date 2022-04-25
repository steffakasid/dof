//go:build release
// +build release

package cmd

import "time"

var buildTime string

func init() {
	dt := time.Now()
	buildTime = dt.Format("2006-01-02T15:04:05")
}
