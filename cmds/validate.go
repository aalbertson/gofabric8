/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cmds

import (
	"strings"

	k8sclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	cmdutil "github.com/GoogleCloudPlatform/kubernetes/pkg/kubectl/cmd/util"
	"github.com/fabric8io/gofabric8/client"
	"github.com/fabric8io/gofabric8/util"
	"github.com/spf13/cobra"
)

type Result string

const (
	Success Result = "✔"
	Failure Result = "✘"
)

type validateFunc func(c *k8sclient.Client, f *cmdutil.Factory) (Result, error)

func NewCmdValidate(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate your Kubernetes or OpenShift environment",
		Long:  `validate your Kubernetes or OpenShift environment`,
		Run: func(cmd *cobra.Command, args []string) {
			c, cfg := client.NewClient(f)
			ns, _, _ := f.DefaultNamespace()
			util.Info("Validating your ")
			util.Success(string(util.TypeOfMaster(c)))
			util.Info(" installation at ")
			util.Success(cfg.Host)
			util.Info(" in namespace ")
			util.Success(ns)
			util.Blank()
			util.Blank()
			printValidationResult("Service account", validateServiceAccount, c, f)
		},
	}

	return cmd
}

func printValidationResult(check string, v validateFunc, c *k8sclient.Client, f *cmdutil.Factory) {
	r, err := v(c, f)
	if err != nil {
		r = Failure
	}
	util.Infof("%s%s", check, strings.Repeat(".", 24-len(check)))
	if r == Failure {
		util.Failuref("%-2s", r)
	} else {
		util.Successf("%-2s", r)
	}
	if err != nil {
		util.Errorf("%v", err)
	}
	util.Blank()
}

func validateServiceAccount(c *k8sclient.Client, f *cmdutil.Factory) (Result, error) {
	ns, _, err := f.DefaultNamespace()
	if err != nil {
		return Failure, err
	}
	sa, err := c.ServiceAccounts(ns).Get("fabric8")
	if sa != nil {
		return Success, err
	}
	return Failure, err
}
