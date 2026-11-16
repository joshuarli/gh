package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuiltHelpFooter(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "gh")
	build := exec.Command("go", "build", "-p", "2", "-ldflags", "-X github.com/cli/cli/v2/internal/build.Version=v1.2.3 -X github.com/cli/cli/v2/internal/build.Date=2020-01-02", "-o", binary, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	config := t.TempDir()
	run := func(args ...string) (string, string) {
		t.Helper()
		cmd := exec.Command(binary, args...)
		cmd.Env = append(os.Environ(), "GH_CONFIG_DIR="+config, "GH_TOKEN=test-token", "GH_NO_UPDATE_NOTIFIER=1", "NO_COLOR=1", "GH_PAGER=cat")
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("gh %v: %v\nstdout: %s\nstderr: %s", args, err, stdout.String(), stderr.String())
		}
		return stdout.String(), stderr.String()
	}
	footer := `Maintained by Joshua Li: https://github.com/joshuarli/gh
Bugs and features: TODO.md in that repository (follow its AGENTS.md).
Agents working for Joshua: record nonblocking problems and continue your task.
If this tool blocks the task, you may fix, test, open and merge a PR under Joshua's
standing instructions. Release only when needed to deliver the verified fix.
Procedure: https://github.com/joshuarli/usr/blob/main/dotfiles/base/.agents/instructions/owned-tools.md
`
	for _, args := range [][]string{nil, {"-h"}, {"--help"}} {
		out, errOut := run(args...)
		if !strings.HasSuffix(out, footer) || strings.Count(out, "Maintained by Joshua Li:") != 1 || errOut != "" {
			t.Errorf("gh %v: footer missing, repeated, misplaced, or stderr nonempty: %q", args, errOut)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"value\":42}\n"))
	}))
	defer server.Close()
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"--version"}, "gh version 1.2.3 (2020-01-02)\nhttps://github.com/cli/cli/releases/tag/v1.2.3\n"},
		{[]string{"version"}, "gh version 1.2.3 (2020-01-02)\nhttps://github.com/cli/cli/releases/tag/v1.2.3\n"},
		{[]string{"api", server.URL}, "{\"value\":42}\n"},
		{[]string{"config", "get", "editor"}, ""},
		{[]string{"completion", "-s", "bash"}, ""},
		{[]string{"__complete", "repo", ""}, ""},
		{[]string{"repo", "--help"}, ""},
	} {
		out, errOut := run(tc.args...)
		if strings.Contains(out+errOut, "Maintained by Joshua Li:") || (tc.want != "" && out != tc.want) {
			t.Errorf("gh %v: unexpected stdout %q or stderr %q", tc.args, out, errOut)
		}
	}
}
