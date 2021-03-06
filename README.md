Package rcfile supplements command-line flags with configuration files.
It is intended as a supplement to the standard `flag` package.

Usage
-----

Define flags using flag.String(), Bool(), Int() etc. like you usually would.
After all flags are defined, call

```
	rcfile.Parse()
	flag.Parse()
```

in that order to parse the user configuration file and command line
arguments into the defined flags. All switches in the configuration file
act as a second set of defaults for all defined flags.

User configuration files for an example program `foobaz` may be in the in the
following locations. Only the first file that exists is parsed.
* ~/.foobazrc
* ${XDG_CONFIG_DIR}/foobazrc
* ~/.config/foobazrc
* %APPDATA%/foobazrc

Configuration file syntax
-------------------------

```
	# Lines starting with hashes are comments
	; as are lines starting with a semicolon
	# Furthermore, empty lines are ignored
	flag1=foo
	flag2  = bar
```

Whitespace is trimmed from both the key and value parts in the config file.
For non-string types, the syntax is identical to the standard `flag` package:

* Integer flags accept `1234`, `0664`, `0x1234` and may be negative.
* Boolean flags may be: `1`, `0`, `t`, `f`, `T`, `F`, `true`, `false`, `TRUE`, `FALSE`, `True`, `False`
* Duration flags accept any input valid for time.ParseDuration.

License
-------
Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.
