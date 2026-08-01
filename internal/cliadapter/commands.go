package cliadapter

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/issam-assiyadi/leftmark"
	"github.com/issam-assiyadi/leftmark/adapter/githook"
	"github.com/issam-assiyadi/leftmark/application"
	"github.com/issam-assiyadi/leftmark/domain"
)

func newService(root string) (*application.Service, error) {
	if root == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		root = wd
	}
	return leftmark.New(root)
}

func printItems(items []domain.Item, asJSON bool, stdout io.Writer) int {
	if asJSON {
		out := make([]itemJSON, len(items))
		for i, item := range items {
			out[i] = toItemJSON(item)
		}
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			return 1
		}
		return 0
	}

	for _, item := range items {
		status := string(item.Status)
		if item.ID != "" && !item.Located {
			status = "missing"
		}
		printf(stdout, "%-9s %-9s %s:%d %s\n", item.Kind, status, item.File, item.Line, item.Text)
	}
	return 0
}

func parseStatuses(csv string) []domain.Status {
	if csv == "" {
		return nil
	}
	var out []domain.Status
	for _, s := range strings.Split(csv, ",") {
		out = append(out, domain.Status(strings.TrimSpace(s)))
	}
	return out
}

func parseKinds(csv string) []domain.Kind {
	if csv == "" {
		return nil
	}
	var out []domain.Kind
	for _, s := range strings.Split(csv, ",") {
		out = append(out, domain.Kind(strings.ToUpper(strings.TrimSpace(s))))
	}
	return out
}

func runScan(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", "", "repo root (default: current directory)")
	jsonOut := fs.Bool("json", false, "print JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	svc, err := newService(*root)
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}

	items, err := svc.Scan()
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}
	return printItems(items, *jsonOut, stdout)
}

func runList(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", "", "repo root (default: current directory)")
	jsonOut := fs.Bool("json", false, "print JSON")
	statusFlag := fs.String("status", "", "comma-separated statuses: open,doing,done")
	kindFlag := fs.String("kind", "", "comma-separated kinds: todo,fixme,note,question")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	svc, err := newService(*root)
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}
	if _, err := svc.Scan(); err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}

	filter := domain.Filter{
		Statuses: parseStatuses(*statusFlag),
		Kinds:    parseKinds(*kindFlag),
	}
	return printItems(svc.List(filter), *jsonOut, stdout)
}

func runPromote(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("promote", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", "", "repo root (default: current directory)")
	jsonOut := fs.Bool("json", false, "print JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	rest := fs.Args()
	if len(rest) != 2 {
		printf(stderr, "usage: leftmark promote <file> <line>\n")
		return 2
	}
	line, err := strconv.Atoi(rest[1])
	if err != nil {
		printf(stderr, "invalid line %q: %v\n", rest[1], err)
		return 2
	}

	svc, err := newService(*root)
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}

	item, err := svc.PromoteToDoing(rest[0], line)
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}
	return printItems([]domain.Item{item}, *jsonOut, stdout)
}

func runResolve(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("resolve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", "", "repo root (default: current directory)")
	jsonOut := fs.Bool("json", false, "print JSON")
	deleteFlag := fs.Bool("delete", false, "delete the tagged comment from the source")
	keepFlag := fs.Bool("keep", false, "keep the tagged comment, just mark done")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	rest := fs.Args()
	if len(rest) != 1 {
		printf(stderr, "usage: leftmark resolve <id> (--delete|--keep)\n")
		return 2
	}
	if *deleteFlag == *keepFlag {
		printf(stderr, "specify exactly one of --delete or --keep\n")
		return 2
	}

	svc, err := newService(*root)
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}

	item, err := svc.ResolveDone(rest[0], *keepFlag)
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}
	return printItems([]domain.Item{item}, *jsonOut, stdout)
}

func runOpen(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("open", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", "", "repo root (default: current directory)")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	rest := fs.Args()
	if len(rest) != 2 {
		printf(stderr, "usage: leftmark open <file> <line>\n")
		return 2
	}
	line, err := strconv.Atoi(rest[1])
	if err != nil {
		printf(stderr, "invalid line %q: %v\n", rest[1], err)
		return 2
	}

	svc, err := newService(*root)
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}

	if err := svc.OpenInEditor(rest[0], line); err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}
	return 0
}

func runReport(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", "", "repo root (default: current directory)")
	jsonOut := fs.Bool("json", false, "print JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	svc, err := newService(*root)
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}
	if _, err := svc.Scan(); err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}

	summary := svc.Report()
	if *jsonOut {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(summary); err != nil {
			printf(stderr, "%v\n", err)
			return 1
		}
		return 0
	}

	printf(stdout, "%d items total\n", summary.Total)
	for status, count := range summary.ByStatus {
		printf(stdout, "  %s: %d\n", status, count)
	}
	if summary.Unlocated > 0 {
		printf(stdout, "%d tracked item(s) not found in the current tree\n", summary.Unlocated)
	}
	return 0
}

func runHook(args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 || args[0] != "install" {
		printf(stderr, "usage: leftmark hook install <pre-commit|pre-push>\n")
		return 2
	}

	hookName := args[1]
	if hookName != "pre-commit" && hookName != "pre-push" {
		printf(stderr, "unknown hook %q, want pre-commit or pre-push\n", hookName)
		return 2
	}

	wd, err := os.Getwd()
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}

	dir := filepath.Join(wd, "scripts", "git-hooks")
	if err := githook.Install(dir, hookName); err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}
	printf(stdout, "installed %s in %s\n", hookName, dir)
	return 0
}
