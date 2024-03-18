package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/pkg/errors"
)

var (
	Version  string = "0.0.1"
	Revision        = func() string { // {{{
		revision := ""
		modified := false
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					//return setting.Value
					revision = setting.Value
					if len(setting.Value) > 7 {
						revision = setting.Value[:7] // 最初の7文字にする
					}
				}
				if setting.Key == "vcs.modified" {
					modified = setting.Value == "true"
				}
			}
		}
		if modified {
			revision = "develop+" + revision
		}
		return revision
	}() // }}}
)

type Args struct {
	Input   []string `arg:"positional"         help:"入力ファイル。"`
	Debug   bool     `arg:"-d,--debug"         help:"デバッグ用。ログが詳細になる。"`
	Version bool     `arg:"-v,--version"       help:"バージョン情報を出力する。"`
}

func (args *Args) Print() {
	log.Printf(`
	Input      : %q
	Debug      : %q
	Version    : %q
`,
		args.Input,
		args.Debug,
		args.Version,
	)
}

type ArgsVersion struct {
}

var parser *arg.Parser

func ShowHelp(post string) {
	buf := new(bytes.Buffer)
	parser.WriteHelp(buf)
	fmt.Printf("%v\n", strings.ReplaceAll(buf.String(), "display this help and exit", "ヘルプを出力する。"))
	if len(post) != 0 {
		fmt.Println(post)
	}
	os.Exit(1)
}
func ShowVersion() {
	if len(Revision) == 0 {
		// go installでビルドされた場合、gitの情報がなくなる。その場合v0.0.0.のように末尾に.がついてしまうのを避ける。
		fmt.Printf("%v version %v\n", GetFileNameWithoutExt(os.Args[0]), Version)
	} else {
		fmt.Printf("%v version %v.%v\n", GetFileNameWithoutExt(os.Args[0]), Version, Revision)
	}
	os.Exit(0)
}

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する
}

func ArgParse(f func(*Args) error) {
	if false {
		if len(os.Args) == 1 {
			// 標準入力から読み取り、標準出力に出力する。
			// ような処理を書くときここに書く。
			os.Exit(0)
		}
	}
	var args *Args = &Args{}
	var err error = nil
	parser, err = arg.NewParser(arg.Config{Program: GetFileNameWithoutExt(os.Args[0]), IgnoreEnv: false}, args)
	if err != nil {
		ShowHelp(fmt.Sprintf("%v", errors.Errorf("%v", err)))
	}
	if err := parser.Parse(os.Args[1:]); err != nil {
		if err.Error() == "help requested by user" {
			ShowHelp("")
		} else if err.Error() == "version requested by user" {
			ShowVersion()
		} else {
			panic(errors.Errorf("%v", err))
		}
	}
	if args.Debug {
		args.Print()
	}
	if args.Version || args.VersionSub != nil {
		//if args.Version {
		ShowVersion()
	}
	// positionalな引数があるとき
	for _, in := range args.Input {
		if err := f(args); err != nil {
			panic(errors.Errorf("%v", err))
		}
		if args.Debug {
			fmt.Println(str)
		}
	}
	if err := f(args); err != nil {
		panic(errors.Errorf("%v", err))
	}
}

func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
func GetFilePathWithoutExt(path string) string {
	return filepath.ToSlash(filepath.Join(filepath.Dir(path), GetFileNameWithoutExt(path)))
}
func ReplaceExt(filePath, from, to string) string {
	ext := filepath.Ext(filePath)
	if len(from) > 0 && ext != from {
		return filePath
	}
	return filePath[:len(filePath)-len(ext)] + to
}
