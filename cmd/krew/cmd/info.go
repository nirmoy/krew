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
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"

	"sigs.k8s.io/krew/pkg/index"
	"sigs.k8s.io/krew/pkg/index/indexscanner"
	"sigs.k8s.io/krew/pkg/installation"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show information about a kubectl plugin",
	Long: `Show information about a kubectl plugin.

This command can be used to print information such as its download URL, last
available version, platform availability and the caveats.

Example:
  kubectl krew info PLUGIN`,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			plugin, err := indexscanner.LoadPluginFileFromFS(paths.IndexPath(), arg)
			if os.IsNotExist(err) {
				return errors.Errorf("plugin %q not found", arg)
			} else if err != nil {
				return errors.Wrap(err, "failed to load plugin manifest")
			}
			printPluginInfo(os.Stdout, plugin)
		}
		return nil
	},
	PreRunE: checkIndex,
	Args:    cobra.MinimumNArgs(1),
}

func printPluginInfo(out io.Writer, plugin index.Plugin) {
	fmt.Fprintf(out, "NAME: %s\n", plugin.Name)
	if platform, ok, err := installation.GetMatchingPlatform(plugin); err == nil && ok {
		if platform.Head != "" {
			fmt.Fprintf(out, "HEAD: %s\n", platform.Head)
		}
		if platform.URI != "" {
			fmt.Fprintf(out, "URI: %s\n", platform.URI)
			fmt.Fprintf(out, "SHA256: %s\n", platform.Sha256)
		}
	}
	if plugin.Spec.Version != "" {
		fmt.Fprintf(out, "VERSION: %s\n", plugin.Spec.Version)
	}
	if plugin.Spec.Homepage != "" {
		fmt.Fprintf(out, "HOMEPAGE: %s\n", plugin.Spec.Homepage)
	}
	if plugin.Spec.Description != "" {
		fmt.Fprintf(out, "DESCRIPTION: \n%s\n", plugin.Spec.Description)
	}
	if plugin.Spec.Caveats != "" {
		fmt.Fprintln(out, prepCaveats(plugin.Spec.Caveats))
	}
}

// prepCaveats converts caveats string to an indented format ready for printing.
// Example:
//
//     CAVEATS:
//     \
//      | This plugin is great, use it with great care.
//      | Also, plugin will require the following programs to run:
//      |  * jq
//      |  * base64
//     /
func prepCaveats(s string) string {
	out := "CAVEATS:\n\\\n"
	s = strings.TrimRightFunc(s, unicode.IsSpace)
	out += regexp.MustCompile("(?m)^").ReplaceAllString(s, " |  ")
	out += "\n/"
	return out
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
