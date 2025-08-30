// Copyright 2025 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build arm64
// +build arm64

package nvidiabfb

import (
	"os"

	"github.com/coreos/ignition/v2/config/v3_6_experimental/types"
	"github.com/coreos/ignition/v2/internal/platform"
	"github.com/coreos/ignition/v2/internal/providers/util"
	"github.com/coreos/ignition/v2/internal/resource"

	"github.com/coreos/vcontext/report"
)

const (
	bootfifoPath = "/sys/bus/platform/devices/MLNXBF04:00/bootfifo"
)

func init() {
	platform.Register(platform.Provider{
		Name:  "nvidiabfb",
		Fetch: fetchConfig,
	})
}

func fetchConfig(f *resource.Fetcher) (cfg types.Config, rpt report.Report, err error) {
	data, err := os.ReadFile(bootfifoPath)
	if os.IsNotExist(err) {
		f.Logger.Info("Nvidia BlueField bootfifo was not found. Ignoring...")
		cfg, rpt, err = util.ParseConfig(f.Logger, []byte{})
		return
	} else if err != nil {
		f.Logger.Err("couldn't read Nvidia BlueField config: %v", err)
		return
	}

	if len(data) == 0 {
		f.Logger.Info("no config found in Nvidia BlueField bootfifo")
		cfg, rpt, err = util.ParseConfig(f.Logger, []byte{})
		return
	}

	cfg, rpt, err = util.ParseConfig(f.Logger, data)
	return
}
