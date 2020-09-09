// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"

	"github.com/deptofdefense/timebucket/pkg/datetime"
)

const (
	TimebucketVersion = "1.0.0"
)

const (
	//
	flagInputFormat = "input-format"
	flagInputColumn = "input-column"
	flagInputValue  = "input-value"
	//
	flagKeyFormat = "key-format"
	//
	flagOutputFormat = "output-format"
	//
	flagLayouts = "layouts"
	flagLimit   = "limit"
	//
	flagSkipErrors = "skip-errors"
	flagTable      = "table"
	//
	flagVersion = "version"
)

func initFlags(flag *pflag.FlagSet) {
	flag.StringP(flagInputFormat, "i", "csv", "input format, one of: csv, jsonl, tags, or tsv")
	flag.StringP(flagInputColumn, "c", "", "input column")
	flag.StringP(flagInputValue, "v", "", "input value as go template")
	flag.StringP(flagKeyFormat, "k", "", "hash key format")
	flag.StringP(flagOutputFormat, "o", "csv", "output format, one of: bson, csv, json, jsonl, properties, tags, tsv, or yaml")
	flag.StringP(flagLayouts, "l", "", "default layouts")
	flag.IntP(flagLimit, "n", serializer.NoLimit, "maximum number of records to process")
	flag.BoolP(flagSkipErrors, "e", false, "skip errors")
	flag.BoolP(flagTable, "t", false, "serialize frequency distribution as table")
	flag.Bool(flagVersion, false, "show version")
}

func initViper(cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()
	err := v.BindPFlags(cmd.Flags())
	if err != nil {
		return v, fmt.Errorf("error binding flag set to viper: %w", err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv() // set environment variables to overwrite config
	return v, nil
}

func checkConfig(v *viper.Viper) error {
	inputFormat := v.GetString(flagInputFormat)
	if len(inputFormat) == 0 {
		return fmt.Errorf("input format is missing")
	}
	if inputFormat != "csv" && inputFormat != "jsonl" && inputFormat != "tags" && inputFormat != "tsv" {
		return fmt.Errorf("input format %q is invalid", inputFormat)
	}
	inputColumn := v.GetString(flagInputColumn)
	inputValue := v.GetString(flagInputValue)
	if len(inputColumn) == 0 && len(inputValue) == 0 {
		return fmt.Errorf("input column and input value are missing, either must be set")
	}
	keyFormat := v.GetString(flagKeyFormat)
	if len(keyFormat) == 0 {
		return fmt.Errorf("key format is missing")
	}
	outputFormat := v.GetString(flagOutputFormat)
	if len(outputFormat) == 0 {
		return fmt.Errorf("output format is missing")
	}
	if inputFormat != "csv" && inputFormat != "jsonl" && inputFormat != "yaml" {
		return fmt.Errorf("input format %q is invalid", inputFormat)
	}
	return nil
}

func main() {
	cmd := &cobra.Command{
		Use:                   `timebucket [flags] -|FILE...`,
		DisableFlagsInUseLine: true,
		Short:                 "timebucket is used to create histograms from temporal data.",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := initViper(cmd)
			if err != nil {
				return fmt.Errorf("error initializing viper: %w", err)
			}

			if v.GetBool(flagVersion) {
				fmt.Println(TimebucketVersion)
				return nil
			}

			if len(args) != 1 {
				return cmd.Usage()
			}

			if errConfig := checkConfig(v); errConfig != nil {
				return errConfig
			}

			inputFormat := v.GetString(flagInputFormat)
			inputColumn := v.GetString(flagInputColumn)
			inputValueTemplate := v.GetString(flagInputValue)
			if len(inputValueTemplate) == 0 {
				inputValueTemplate = fmt.Sprintf("{{.%s}}", inputColumn)
			}
			inputValue, err := template.New("main").Parse(inputValueTemplate)
			if err != nil {
				return fmt.Errorf("error parsing input template: %w", err)
			}
			keyFormat := v.GetString(flagKeyFormat)
			outputFormat := v.GetString(flagOutputFormat)
			limit := v.GetInt(flagLimit)
			layouts := datetime.DefaultLayouts
			if str := v.GetString(flagLayouts); len(str) > 0 {
				layouts = strings.Split(str, ",")
			}
			skipErrors := v.GetBool(flagSkipErrors)
			table := v.GetBool(flagTable)

			var reader io.Reader
			if args[0] == "-" {
				reader = bufio.NewReader(os.Stdin)
			} else {
				r, errOpen := os.Open(args[0])
				if errOpen != nil {
					return fmt.Errorf("error opening input file %q: %w", args[0], errOpen)
				}
				reader = r
			}

			it, err := iterator.NewIterator(&iterator.NewIteratorInput{
				Reader:            reader,
				Format:            inputFormat,
				Header:            nil,
				ScannerBufferSize: 4096,
				SkipLines:         serializer.NoSkip,
				SkipBlanks:        true,
				SkipComments:      true,
				Comment:           serializer.NoComment,
				Trim:              false,
				LazyQuotes:        true,
				Limit:             limit,
				KeyValueSeparator: "=",
				LineSeparator:     "\n",
				DropCR:            true,
				Type:              reflect.TypeOf(map[string]interface{}{}),
			})
			if err != nil {
				return fmt.Errorf("error creating iterator: %w", err)
			}

			freqdist := map[string]int{}

			for {

				// next object
				obj, errNext := it.Next()
				if errNext != nil {
					if errNext == io.EOF {
						break
					}
					return fmt.Errorf("error reading record: %w", errNext)
				}

				// calculate value
				buf := new(bytes.Buffer)
				errExecute := inputValue.Execute(buf, obj)
				if errExecute != nil {
					if skipErrors {
						continue
					} else {
						return fmt.Errorf("error creating input value from record: %w", errExecute)
					}
				}

				// parse value
				v, errParse := datetime.Parse(buf.String(), layouts)
				if errParse != nil {
					if skipErrors {
						continue
					} else {
						return fmt.Errorf("error parsing date from record: %w", errParse)
					}
				}

				// add to histogram
				key := v.Format(keyFormat)
				if _, ok := freqdist[key]; !ok {
					freqdist[key] = 0
				}
				freqdist[key] += 1
			}

			// serialize frequency distribution as an object
			if (!table) && (outputFormat == "bson" || outputFormat == "json" || outputFormat == "properties" || outputFormat == "yaml") {
				outputBytes, errNew := serializer.New(outputFormat).
					LineSeparator("\n").
					KeyValueSeparator("=").
					Serialize(freqdist)
				if errNew != nil {
					return fmt.Errorf("error serializing output: %w", errNew)
				}
				_, errWrite := os.Stdout.Write(outputBytes)
				if errWrite != nil {
					return fmt.Errorf("error writing output: %w", errWrite)
				}
				return nil
			}

			// collect keys
			keys := make([]string, 0, len(freqdist))
			for k := range freqdist {
				keys = append(keys, k)
			}

			// sort keys
			sort.Strings(keys)

			// create rows
			outputObject := make([]interface{}, 0, len(keys))
			for _, k := range keys {
				outputObject = append(outputObject, map[string]interface{}{
					"key":   k,
					"count": freqdist[k],
				})
			}

			// serialize frequency distribution as a table
			outputBytes, err := serializer.New(outputFormat).
				Header([]interface{}{"key", "count"}).
				LineSeparator("\n").
				KeyValueSeparator("=").
				Serialize(outputObject)
			if err != nil {
				return fmt.Errorf("error serializing output: %w", err)
			}
			_, err = os.Stdout.Write(outputBytes)
			if err != nil {
				return fmt.Errorf("error writing output: %w", err)
			}
			return nil
		},
	}
	initFlags(cmd.Flags())

	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "timebucket: "+err.Error())
		_, _ = fmt.Fprintln(os.Stderr, "Try timebucket --help for more information.")
		os.Exit(1)
	}
}
