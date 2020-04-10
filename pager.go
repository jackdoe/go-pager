package pager

import (
	"io"
	"log"
	"os"
	"os/exec"
)

// find pager path using PAGER env variable or attempt the array of arguments
// e.g. GetPagerPath("less","more","cat")
func getPagerPath(try ...string) string {
	p := os.Getenv("PAGER")
	if p != "" {
		if p == "NOPAGER" {
			return ""
		}

		exe, err := exec.LookPath(p)
		if err != nil {
			log.Fatal(err)
		}
		return exe
	}

	for _, x := range try {
		exe, err := exec.LookPath(x)
		if err == nil {
			return exe
		}
	}

	return ""
}

// Create new pager executing a command based on $PAGER env var or
// array of executables
// Example:
//
//	p, close := pager.Pager("less", "more","cat")
//	defer close()
//
//	p.Write("hello world")
//
// Will try to find $PAGER,less,more,cat in path, first one it finds it will
// pipe the output written to the returned Writer, if nothing is found
// it will return os.Stdout, if $PAGER=NOPAGER it will return
// os.Stdout, if $PAGER is specified but cant be found in path it will
// panic
func Pager(try ...string) (io.Writer, func()) {
	p := getPagerPath(try...)
	if p != "" {
		cmd := exec.Command(p)
		r, w := io.Pipe()
		cmd.Stdin = r
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		c := make(chan struct{})
		go func() {
			defer close(c)
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
			os.Exit(0)
		}()

		return w, func() {
			w.Close()
			<-c
		}
	}
	return os.Stdout, func() {}
}
