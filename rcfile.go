package rcfile

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)

func Parse() {
	f := openRCFile()
	if f == nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	i := 0
	for {
		l, err := r.ReadString('\n')
		if err != nil {
			break
		}
		i++
		l = strings.Trim(l, " \t\v\r\n")

		if len(l) == 0 {
			continue
		}
		if l[0] == '#' || l[0] == ';' {
			continue
		}

		var i int
		var c rune
		for i, c = range l {
			if c == '=' {
				break
			}
		}

		if c != '=' {
			log.Fatalf("Syntax error in config file on line %d: Expected key=value pair; got: '%s'", i, l)
		}

		k := strings.Trim(l[0:i], " \t\v\r\n")
		v := strings.TrimLeft(l[i+1:], " \t\v\r\n")

		er := flag.Set(k, v)
		if er != nil {
			log.Fatal(er)
		}
	}
}

func openRCFile() io.ReadCloser {
	var err error
	if u, err := user.Current(); err == nil {
		rcf := path.Join(u.HomeDir, "."+path.Base(os.Args[0])+"rc")
		r, err := os.Open(rcf)
		if err != nil {
			log.Print(err)
		} else {
			return r
		}
	}
	log.Print(err)
	return nil
}
