// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckPlatformCompat(t *testing.T) {
	t.Skip("there is no failure condition on this test and it's flaky, " +
		"failing due to the default 10min timeout. " +
		"See https://github.com/elastic/elastic-agent/issues/3964")
	if !(runtime.GOARCH == "amd64" && (isLinux() ||
		isWindows())) {
		t.Skip("Test not support on current platform")
	}

	// compile test helper
	tmp := t.TempDir()
	helper := filepath.Join(tmp, "helper")

	cmd := exec.Command("go", "test", "-c", "-o", helper)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GOARCH=386")
	require.NoError(t, cmd.Run(), "failed to compile test helper")
	t.Logf("compiled test binary %q", helper)

	// run test helper
	cmd = exec.Command(helper, "-test.v", "-test.run", "TestHelper")
	cmd.Env = []string{"GO_USE_HELPER=1"}
	t.Logf("running %q", cmd.Args)
	output, err := cmd.Output()
	if err != nil {
		t.Logf("32bit binary tester failed.\n Err: %v\nOutput: %s",
			err, output)
	}
}

func isLinux() bool {
	return runtime.GOOS == "linux"
}

func TestHelper(t *testing.T) {
	if os.Getenv("GO_USE_HELPER") != "1" {
		t.Log("ignore helper")
		return
	}

	err := CheckNativePlatformCompat()
	if err.Error() != "trying to run 32Bit binary on 64Bit system" {
		t.Error("expected the native platform check to fail")
	}
}
