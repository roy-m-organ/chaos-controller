// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2020 Datadog, Inc.

package main

import (
	"github.com/DataDog/chaos-controller/api/v1beta1"
	"github.com/DataDog/chaos-controller/container"
	"github.com/DataDog/chaos-controller/injector"
	"github.com/spf13/cobra"
)

var cpuPressureInjectCmd = &cobra.Command{
	Use:   "inject",
	Short: "Inject a CPU pressure on the actual node",
	Run: func(cmd *cobra.Command, args []string) {
		uid, _ := cmd.Flags().GetString("uid")
		containerid, _ := cmd.Flags().GetString("container-id")

		// prepare container
		ctn, err := container.New(containerid)
		if err != nil {
			log.Fatalw("can't create container object", "error", err)
		}

		// prepare spec
		spec := v1beta1.CPUPressureSpec{}

		// inject
		i := injector.NewCPUPressureInjector(uid, spec, ctn, log, ms)
		i.Inject()
	},
}