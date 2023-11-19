package main

import (
	"bytes"
	"context"
	_ "embed"
	json2 "encoding/json"
	"flag"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
	"strings"
	"text/template"
)

type ObjectWithExists[TIdent, TObj any] struct {
	Identifier TIdent
	Object     TObj
	Exists     bool
}

type EnvFlag struct {
	Dest   *string
	EnvVar string
	Option string
}

type SubcommandTemplate[TArg, TOut any] struct {
	// --------------------
	// Private fields
	flagSet        *flag.FlagSet
	envFlags       []EnvFlag
	showSubCmdHelp bool
	showSubCmdJson bool

	// ---------------------
	// public fields
	Command  []string
	Template *template.Template

	Arg TArg

	FlagSetInit    func(arg *TArg, fs *flag.FlagSet)
	EnvFlagSetInit func(arg *TArg) []EnvFlag

	Validator func(arg *TArg) error
	Executor  func(context.Context, v3client.Client, TArg) (TOut, error)
}

type SubcommandFinder struct {
	Command  []string
	Executor ExecutorFunc
}

func (st *SubcommandTemplate[TArg, TOut]) ExecuteCLI(ctx context.Context, client v3client.Client, subCmd []string) int {
	return st.Execute(ctx, client, subCmd)
}

func (st *SubcommandFinder) Matches(args []string) int {
	if len(args) >= len(st.Command) {
		for i := 0; i < len(st.Command); i++ {
			if strings.ToLower(args[i]) != st.Command[i] {
				return 0
			}
		}

		return len(st.Command)
	} else {
		return 0
	}
}

func (st *SubcommandTemplate[TArg, TOut]) Finder() *SubcommandFinder {
	return &SubcommandFinder{
		Command:  st.Command,
		Executor: st.ExecuteCLI,
	}
}

func mustTemplate(str string) *template.Template {
	if t, err := template.New("templ").Parse(str); err != nil {
		panic(err.Error())
	} else {
		return t
	}
}

func executeTemplate(st *template.Template, obj any) (string, int) {
	sb := bytes.Buffer{}
	if templErr := st.Execute(&sb, obj); templErr != nil {
		return fmt.Sprintf("template error %s", templErr.Error()), 3
	} else {
		return sb.String(), 0
	}
}

func (st *SubcommandTemplate[TArg, TOut]) parseCommand(args []string) error {
	st.flagSet = flag.NewFlagSet(strings.Join(st.Command, " "), flag.ContinueOnError)
	st.flagSet.BoolVar(&st.showSubCmdHelp, "help", false, "Show sub-command help")
	st.flagSet.BoolVar(&st.showSubCmdJson, outputJsonOps, false, "Render output as json")

	if st.FlagSetInit != nil {
		st.FlagSetInit(&st.Arg, st.flagSet)
	}

	if parseErr := st.flagSet.Parse(args); parseErr != nil {
		return parseErr
	}

	// Load data from the environment set
	if st.EnvFlagSetInit != nil {
		setters := st.EnvFlagSetInit(&st.Arg)
		for _, s := range setters {
			if len(*s.Dest) == 0 {
				if envVal := os.Getenv(s.EnvVar); len(envVal) > 0 {
					*s.Dest = envVal
				}
			}
		}
	}

	return nil
}

//go:embed templates/subcmd_env_flag_set.tmpl
var subcmdEnvFlagSetTemplate string

func (st *SubcommandTemplate[TArg, TOut]) Execute(ctx context.Context, cl v3client.Client, args []string) int {
	cmdParseErr := st.parseCommand(args[len(st.Command):])
	if cmdParseErr != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error parsing CLI command: %s", cmdParseErr.Error()))
		return 1
	}

	if st.showSubCmdHelp {
		st.flagSet.PrintDefaults()

		if st.EnvFlagSetInit != nil {
			if envOpts := st.EnvFlagSetInit(&st.Arg); len(envOpts) > 0 {
				fsTempl := mustTemplate(subcmdEnvFlagSetTemplate)
				fsTempl.Execute(os.Stdout, envOpts)
			}
		}

		return 0
	}

	if st.Validator != nil {
		if validateErr := st.Validator(&st.Arg); validateErr != nil {
			os.Stderr.WriteString(fmt.Sprintf("Input is not valid for this command: %s\n", validateErr.Error()))
			return 1
		}
	}

	if tOut, execErr := st.Executor(ctx, cl, st.Arg); execErr != nil {
		os.Stderr.WriteString(fmt.Sprintf("Command execution has failed: %s\n", execErr.Error()))
		return 2
	} else {
		if globalOptOutputJson || st.showSubCmdJson {
			sb := strings.Builder{}
			encoder := json2.NewEncoder(&sb)
			encoder.SetIndent("", "  ")
			if encErr := jsonEncoder.Encode(tOut); encErr != nil {
				os.Stderr.WriteString(fmt.Sprintf("Could not produce JSON: %s\n", encErr.Error()))
				return 23
			} else {
				fmt.Println(sb.String())
				return 0
			}

		} else {
			output, rv := executeTemplate(st.Template, tOut)
			fmt.Println(output)
			return rv
		}
	}
}
