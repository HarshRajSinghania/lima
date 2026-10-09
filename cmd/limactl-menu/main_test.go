// SPDX-FileCopyrightText: Copyright The Lima Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"strings"
	"testing"

	"github.com/lima-vm/lima/v2/pkg/limatype"
)

func TestFormatIECUsesBinaryUnits(t *testing.T) {
	const gib = 1024 * 1024 * 1024
	if got := formatIEC(4 * gib); got != "4GiB" {
		t.Fatalf("4GiB rendered as %q, want 4GiB", got)
	}
	if got := formatIEC(512 * 1024 * 1024); got != "512MiB" {
		t.Fatalf("512MiB rendered as %q", got)
	}
	if got := formatIEC(1536 * 1024 * 1024); got != "1.5GiB" {
		t.Fatalf("1.5GiB rendered as %q", got)
	}
	if got := formatIEC(0); got != "0B" {
		t.Fatalf("0 rendered as %q", got)
	}
	if got := formatIEC(-1); got != "-" {
		t.Fatalf("negative rendered as %q", got)
	}
}

func TestFmtInstSeparatesCPUAndMemory(t *testing.T) {
	in := &limatype.Instance{
		Name:         "default",
		Status:       "Running",
		Arch:         "aarch64",
		CPUs:         4,
		Memory:       4 * 1024 * 1024 * 1024,
		SSHAddress:   "127.0.0.1",
		SSHLocalPort: 58684,
	}
	line := stripTviewTags(fmtInst(in))
	if strings.Contains(line, "4.3GB") || strings.Contains(line, "4CPU/") {
		t.Fatalf("memory still combined or decimal: %s", line)
	}
	if !strings.Contains(line, "4GiB") {
		t.Fatalf("missing 4GiB in %s", line)
	}
	header := instHeader()
	memAt := strings.Index(header, "MEMORY")
	if memAt < 0 {
		t.Fatal("header missing MEMORY")
	}
	if got := strings.TrimSpace(line[memAt : memAt+colMem]); got != "4GiB" {
		t.Fatalf("MEMORY column = %q in %q (header %q)", got, line, header)
	}
	cpuAt := strings.Index(header, "CPU")
	if got := strings.TrimSpace(line[cpuAt : cpuAt+colCPU]); got != "4" {
		t.Fatalf("CPU column = %q in %q", got, line)
	}
}

func stripTviewTags(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '[' {
			b.WriteByte(s[i])
			continue
		}
		end := strings.IndexByte(s[i+1:], ']')
		if end < 0 {
			b.WriteByte(s[i])
			continue
		}
		i += end + 1
	}
	return b.String()
}
