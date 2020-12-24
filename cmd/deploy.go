/*
Copyright © 2020 Niall Grant <ngfgrant@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"ideme/do"
	"strings"
)

var Image string

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your infrastructure or application.",
}

var deployInfrastructureCmd = &cobra.Command{
	Use:   "infrastructure",
	Short: "Deploy your Infrastructure",
	Run: func(cmd *cobra.Command, args []string) {
		api := do.ConfigureApi()
		do.CreateInfrastructure(api)
	},
}

var deployAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Deploy your application",
	Run: func(cmd *cobra.Command, args []string) {
		api := do.ConfigureApi()
		infra := do.GetInfrastructure(api)
		var appImage string = Image

		if Image == "" {
			appImage = viper.GetString("application.image")
		}

		image := strings.ReplaceAll(appImage, "/", "\\/")
		do.CreateApplication(api, infra, image)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.AddCommand(deployInfrastructureCmd)
	deployCmd.AddCommand(deployAppCmd)
	deployAppCmd.Flags().StringVarP(&Image, "image", "i", "", "Docker image to use.")
}
