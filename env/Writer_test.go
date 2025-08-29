package env

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestWriterSetUnset(t *testing.T) {
	w := &Writer{}
	w.Set("GREETING", "Hello World")
	w.Set("DEBUG", "true")
	w.Unset("DEBUG")
	script := w.String()
	if !strings.Contains(script, "hex2str") {
		t.Errorf("expected script to contain hex2str function, got:\n%s", script)
	} else if !strings.Contains(script, "export GREETING=$(hex2str") {
		t.Errorf("expected export GREETING line, got:\n%s", script)
	} else if !strings.Contains(script, "unset DEBUG") {
		t.Errorf("expected unset DEBUG line, got:\n%s", script)
	}
	if t.Failed() {
		t.FailNow()
	}
	got := runScript(t, script+`echo "$GREETING"`)
	if want := "Hello World"; got != want {
		t.Errorf("expected GREETING=%q, got %q", want, got)
	}
	got = runScript(t, script+`
		if [ -z "$DEBUG" ]; then
			echo "DEBUG is unset"
		else
			echo "DEBUG is set"
		fi
	`)
	if want := "DEBUG is unset"; got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
func TestWriterSetAndUnsetBranches(t *testing.T) {
	w := &Writer{}
	w.Set("FOO", "bar")
	w.Set("BAZ", "qux")
	w.Unset("FOO")
	w.Reset()
	w.Unset("DEBUG")
	script := w.String()
	if !strings.Contains(script, "hex2str") {
		t.Errorf("expected embedded hex2str function in script")
	} else if !strings.Contains(script, "unset DEBUG") {
		t.Errorf("expected unset DEBUG after Reset, got:\n%s", script)
	}
}
func TestWriterBytesAndString(t *testing.T) {
	w := &Writer{}
	w.Set("FOO", "bar")
	if string(w.Bytes()) != w.String() {
		t.Errorf("Bytes and String mismatch")
	}
}
func TestWriterWriteTo(t *testing.T) {
	w := &Writer{}
	w.Set("FOO", "bar")
	var buf bytes.Buffer
	n, err := w.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	} else if int64(buf.Len()) != n {
		t.Errorf("expected n=%d, got %d", buf.Len(), n)
	} else if !strings.Contains(buf.String(), "FOO") {
		t.Errorf("WriteTo output missing FOO, got:\n%s", buf.String())
	}
}
func TestWriterWrite(t *testing.T) {
	w := &Writer{}
	_, _ = w.Write([]byte("hello"))
	if !strings.Contains(w.String(), "hello") {
		t.Errorf("expected 'hello' in buffer, got:\n%s", w.String())
	}
}
func TestWriterIntegration(t *testing.T) {
	w := &Writer{}
	w.Set("GREETING", "Hello World")
	w.Unset("DEBUG")
	script := w.String()
	got := runScript(t, script+`
		echo "$GREETING"
	`)
	if got != "Hello World" {
		t.Errorf("expected GREETING=%q, got %q", "Hello World", got)
	}
	got = runScript(t, script+`
		if [ -z "$DEBUG" ]; then
			echo "DEBUG is unset"
		else
			echo "DEBUG is set"
		fi
	`)
	if got != "DEBUG is unset" {
		t.Errorf("expected DEBUG unset, got %q", got)
	}
}
func runScript(t *testing.T, script string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		possible := []string{
			`C:\Program Files\Git\bin\bash.exe`,
			`C:\Program Files (x86)\Git\bin\bash.exe`,
		}
		for _, p := range possible {
			if _, err := os.Stat(p); err == nil {
				return runWith(t, p, script)
			}
		}
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH; skipping integration test")
	}
	return runWith(t, bash, script)
}
func runWith(t *testing.T, bash string, script string) string {
	t.Helper()
	cmd := exec.Command(bash, "-c", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run script (%s): %v\nOutput:\n%s", bash, err, out)
	}
	return strings.TrimSpace(string(out))
}
