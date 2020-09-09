# timebucket

## Description

`timebucket` is used to create histograms from temporal data.  Timebucket uses streaming to read data, so that memory consumption is low.  Timebucket can read data formatted as `csv`, `jsonl`, `tags`, or `tsv`.  Timebucket can format the histogram as `bson`, `csv`, `json`, `jsonl`, `properties`, `tags`, `tsv`, or `yaml`.

## Usage

```text
timebucket is used to create histograms from temporal data.

Usage:
  timebucket [flags] -|FILE...

Flags:
  -h, --help                   help for timebucket
  -c, --input-column string    input column
  -i, --input-format string    input format, one of: csv, jsonl, tags, or tsv (default "csv")
  -v, --input-value string     input value as go template
  -k, --key-format string      hash key format
  -l, --layouts string         default layouts
  -n, --limit int              maximum number of records to process (default -1)
  -o, --output-format string   output format, one of: bson, csv, json, jsonl, properties, tags, tsv, or yaml (default "csv")
  -e, --skip-errors            skip errors
  -t, --table                  serialize frequency distribution as table
      --version                show version
```

### Selecting

In order to put each record in a bucket, a single date time string must be selected and parsed.  You can use either the `--input-column` or `--input-value` flags to select the string value to parse for each entry.  The `--input-column` (`-c`) flag is used to select a single field.  The `--input-value` (`-v`) flag is used to combine multiple fields using the go template language.  Examples are shown below.

### Parsing

`timebucket` uses the native Go string format for parsing times, which can be confusing.  The placeholders in a time format are actual numbers.  For example, `15` specifies the hour.  Yes, actually use `15` and not `HH`!  See the [time](https://pkg.go.dev/time?tab=doc#pkg-constants) package for more information.  In Timebucket, The input string value for each object is parsed by iterating through a list of formats.  The first format that parses without error is used.

```go
var DefaultLayouts = []string{
	// timestamps
	"1/2/06 15:04:05 PM (MST)",
	"1/2/06 15:04:05 PM",
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"02 Jan 2006T15:04:05.999999999",
	// date only
	"2006",
	"2006-01-02",
	// time only
	"15:04:05.999999999",
	"15:04:05",
}
```

To parse the input value using a custom set of layouts use the `--layouts` flag with a comma separated list of layouts.

### Docker

If you wish to build a Docker container for running `timebucket`, then use the following line.

```shell
make docker_build
```

To run the `timebucket` Docker container, using the following usage line.

```shell
docker run -it --rm=true -v ${PWD}:${PWD} -w ${PWD} timebucket:latest [flags] FILE
```

## Examples

The below commands reads in the data from `input.csv` and counts by hour.

```shell
timebucket --input-format csv --input-column Time --key-format '15' --output-format csv input.csv > hours.csv
```

You can also read data from stdin by using the "-" path.

```shell
cat input.csv.gz | gunzip | timebucket --input-format csv --input-column Time --key-format '15' --output-format csv - > hours.csv
```

The below commands reads in the data from `input.csv`, parses using a custom layout, and counts by year.

```shell
timebucket --input-format csv --input-column Year --layouts '2006' --key-format '2006' --output-format csv input.csv > hours.csv
```

You can use the `--input-value` (`-v`) flag to combine multiple input value and parse.  For example, the below command counts the number of hurricanes by month.

```shell
timebucket -v '{{.DayOfMonth}} {{.Year}}' --layouts 'January 2 2006' -k 'January' input.csv
```

If you select an output format of `bson`, `json`, `properties`, or `yaml`, then you can use the `--table` (`-t`) flag to serialize the frequency distribution as a table rather than an object.

## Building

**timebucket** is written in pure Go, so the only dependency needed to compile the program is [Go](https://golang.org/).  Go can be downloaded from <https://golang.org/dl/>.

This project uses [direnv](https://direnv.net/) to manage environment variables and automatically adding the `bin` and `scripts` folder to the path.  Install direnv and hook it into your shell.  The use of `direnv` is optional as you can always call iceberg directly with `bin/iceberg`.

If using `macOS`, follow the `macOS` instructions below.

To build a binary for your local operating system you can use `make bin/timebucket`.  To build for a release, you can use `make build_release`.  Additionally, you can call `go build` directly to support specific use cases.

### macOS

You can install `go` on macOS using homebrew with `brew install go`.

To install `direnv` on `macOS` use `brew install direnv`.  If using bash, then add `eval \"$(direnv hook bash)\"` to the `~/.bash_profile` file .  If using zsh, then add `eval \"$(direnv hook zsh)\"` to the `~/.zshrc` file.

## Contributing

We'd love to have your contributions!  Please see [CONTRIBUTING.md](CONTRIBUTING.md) for more info.

## Security

Please see [SECURITY.md](SECURITY.md) for more info.

## License

This project constitutes a work of the United States Government and is not subject to domestic copyright protection under 17 USC § 105.  However, because the project utilizes code licensed from contributors and other third parties, it therefore is licensed under the MIT License.  See LICENSE file for more information.
