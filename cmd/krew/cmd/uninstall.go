// Copyright 2019 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"sigs.k8s.io/krew/pkg/installation"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall plugins",
	Long: `Uninstall one or more plugins.

Example:
  kubectl krew uninstall NAME [NAME...]

Remarks:
  Failure to uninstall a plugin will result in an error and exit immediately.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, name := range args {
			glog.V(4).Infof("Going to uninstall plugin %s\n", name)
			if err := installation.Uninstall(paths, name); err != nil {
				return errors.Wrapf(err, "failed to uninstall plugin %s", name)
			}
			fmt.Fprintf(os.Stderr, "Uninstalled plugin %s\n", name)
		}
		return nil
	},
	PreRunE: checkIndex,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"remove"},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
