package rcfile

/*
	Package rcfile supplements command-line flags with configuration files.
	It is intended as a supplement to the standard `flag` package.

	Usage:

	Define flags using flag.String(), Bool(), Int() etc. like you usually would.
	After all flags are defined, call
		rcfile.Parse()
		flag.Parse()
	in that order to parse the user configuration file and command line
	arguments into the defined flags. All switches in the configuration file
	act as a second set of defaults for all defined flags.

	User configuration files for an example program `foobaz` may be in the in the
	following locations. Only the first file that exists is parsed.
		~/.foobazrc
		${XDG_CONFIG_DIR}/foobazrc
		~/.config/foobazrc
		%APPDATA%/foobazrc

	Configuration file syntax:
		# Lines starting with hashes are comments
		; as are lines starting with a semicolon
		# Furthermore, empty lines are ignored
		flag1=foo
		flag2  = bar

	Whitespace is trimmed from both the key and value parts in the config file.
	For non-string types, the syntax is identical to the standard `flag` package:

	Integer flags accept 1234, 0664, 0x1234 and may be negative.
	Boolean flags may be:
		1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False
	Duration flags accept any input valid for time.ParseDuration.
*/

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
	rc := path.Base(os.Args[0]) + "rc"

	user, uexists := user.Current()
	if uexists == nil {
		rcf := path.Join(user.HomeDir, "."+rc)
		r, err := os.Open(rcf)
		if err == nil {
			return r
		}
	}

	xdg_config := os.Getenv("XDG_CONFIG_DIR")
	if xdg_config != "" {
		r, err := os.Open(path.Join(xdg_config, rc))
		if err == nil {
			return r
		}
	}

	if uexists == nil {
		r, err := os.Open(path.Join(user.HomeDir, ".config", rc))
		if err == nil {
			return r
		}
	}

	appdata := os.Getenv("APPDATA")
	if appdata != "" {
		r, err := os.Open(path.Join(appdata, rc))
		if err == nil {
			return r
		}
	}

	return nil
}
