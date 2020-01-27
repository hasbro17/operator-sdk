// Copyright 2018 The Operator-SDK Authors
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

package generate

import (
	"github.com/spf13/cobra"

	"github.com/operator-framework/operator-sdk/cmd/operator-sdk/generate"
)

// Directory names in Kubebuilder's project layout
const (
	configDir = "config"
	apisDir   = "api"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate <generator>",
		Short: "Invokes a generate command",
		Long: `The operator-sdk generate command invokes a command to perform certain code-
and manifest-generation tasks.`,
	}
	cmd.AddCommand(
		// Reuse the existing `generate csv` cmd after configuring it for the kubebuilder layout
		generate.NewGenerateCSVCmd(&generate.CSVCmdConfig{
			IncludePaths:  []string{configDir},
			DeployDirPath: configDir,
			APIsDirPath:   apisDir,
		}),
	)
	return cmd
}
