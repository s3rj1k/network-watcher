package main

import (
	"fmt"
	"os"
)

const (
	defaultAbbRevisionNum = 8
)

func printInfo() {
	vcsRev, vcsDate := getVCSInfo(defaultAbbRevisionNum)
	if vcsDate.Unix() == 0 {
		_, _ = fmt.Fprintf(os.Stdout, "VCS_REV: %s\n",
			vcsRev,
		)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "VCS_REV: %s | VCS_DATE: %v\n",
			vcsRev, vcsDate,
		)
	}
}
