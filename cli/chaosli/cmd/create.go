// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2021 Datadog, Inc.

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/DataDog/chaos-controller/api/v1beta1"
	"github.com/DataDog/chaos-controller/types"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a disruption.",
	Long:  `creates a disruption given input from the user.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(create())
	},
}

const intro = `Hello! This tool will walk you through creating a disruption. Please reply to the prompts, and use ctrl+c to end.
The generated disruption will have "dryRun:true" set for safety, which means you can safely apply it without injecting any failure.`

func create() (v1beta1.DisruptionSpec, error) {
	fmt.Println(intro)

	spec := v1beta1.DisruptionSpec{}
	spec.DryRun = true

	err := promptForKind(&spec)

	if err != nil {
		return spec, err
	}

	spec.Level = getLevel()
	spec.Selector = getSelectors()
	spec.Count = getCount()

	if spec.Level == types.DisruptionLevelPod {
		spec.Containers = getContainers()
	}

	return spec, nil
}

func indexOfString(slice []string, indexed string) int {
	for i, item := range slice {
		if item == indexed {
			return i
		}
	}

	return -1
}

func promptForKind(spec *v1beta1.DisruptionSpec) error {
	initial := "Let's begin by choosing the type of disruption to apply! Which disruption kind would you like to add?"
	followUp := "Would you like to add another disruption kind? It's not necessary, most disruptions involve only one kind. Ctrl+C to finish adding kinds."
	kinds := []string{"dns", "network", "cpu", "disk", "node failure"}
	helpText := `The DNS disruption allows for overriding the A or CNAME records returned by DNS queries.
The Network disruption allows for injecting a variety of different network issues into your target.
The CPU and Disk Pressure disruptions apply cpu or IO pressure to your target, respectively.
Tne Node Failure disruption can either shutdown or restart the targeted node, or the node hosting the targeted pod.

Select one for more information on it.`

	query := initial

	for {
		response, err := selectInput(query, kinds, helpText)
		if err != nil && (query != followUp || err != terminal.InterruptErr) {
			// An initial ctrl+c means we should abort the entire interaction, but
			// once we're on follow-up, that merely means the user is done
			return err
		} else if query == followUp && err == terminal.InterruptErr {
			return nil
		}

		switch response {
		case "dns":
			spec.DNS = getDNS()
			if spec.DNS == nil {
				continue
			}
		case "network":
			spec.Network = getNetwork()
			if spec.Network == nil {
				continue
			}
		case "cpu":
			spec.CPUPressure = getCPUPressure()
			if spec.CPUPressure == nil {
				continue
			}
		case "disk":
			spec.DiskPressure = getDiskPressure()
			if spec.DiskPressure == nil {
				continue
			}
		case "node failure":
			spec.NodeFailure = getNodeFailure()
			if spec.NodeFailure == nil {
				continue
			}
		}

		query = followUp
		i := indexOfString(kinds, response)
		kinds = append(kinds[:i], kinds[i+1:]...)
	}
}

func confirmKind(kind string, helpText string) bool {
	return confirmOption(fmt.Sprintf("Would you like to include the disruption %s?", kind), helpText)
}

func confirmOption(query string, helpText string) bool {
	var result bool

	prompt := &survey.Confirm{
		Message: query,
		Help:    helpText,
	}

	err := survey.AskOne(prompt, &result)

	if err != nil {
		fmt.Printf("confirmOption failed: %v", err)
	}

	return result
}

func getInput(query string, helpText string) string {
	var result string

	prompt := &survey.Input{
		Message: query,
		Help:    helpText,
	}

	err := survey.AskOne(prompt, &result)

	if err != nil {
		fmt.Printf("getInput failed: %v", err)
	}

	return result
}

func selectInput(query string, inputs []string, helpText string) (string, error) {
	var result string

	prompt := &survey.Select{
		Message: query,
		Options: inputs,
		Help:    helpText,
	}

	err := survey.AskOne(prompt, &result)

	return result, err
}

func getSliceInput(query string, helpText string) []string {
	var results string

	prompt := &survey.Multiline{
		Message: query,
		Help:    helpText,
	}

	err := survey.AskOne(prompt, &results)

	if err != nil {
		fmt.Printf("getSliceInput failed: %v\n", err)
	}

	return strings.Split(results, "\n")
}

func getDNS() v1beta1.DNSDisruptionSpec {
	if !confirmKind("DNS Disruption", "Overrides DNS resolution for specified hostnames with a MitM DNS attack. All other DNS requests will use the target's normal DNS resolver.") {
		return nil
	}

	getHostRecordPair := func() v1beta1.HostRecordPair {
		hrPair := v1beta1.HostRecordPair{}

		hrPair.Hostname = getInput("Specify a hostname to target",
			"When your target makes a DNS request for this hostname; the disruption will make sure the value you specify is returned, rather than the real record.")
		hrPair.Record.Type, _ = selectInput("the type of DNS record to inject",
			[]string{"A", "CNAME"},
			"We only support these two types of DNS requests for now. An A record request gets back an IP for a hostname, while a CNAME request maps an alias domain name to the canonical name.")
		helpText := "We're specifying an A record, so the value should be an IP address. You can specify multiple IP addresses, if desired. Simply delimit them with commas, no whitespace! The disruption will round-robin between the options."

		if hrPair.Record.Type == "CNAME" {
			helpText = "We're specifying a CNAME record, so the value should be a hostname to redirect to."
		}

		hrPair.Record.Value = getInput("What value would you like to inject into this DNS record?", helpText)

		return hrPair
	}

	fmt.Println("Let's specify a DNS record to inject!")

	spec := v1beta1.DNSDisruptionSpec{}
	spec = append(spec, getHostRecordPair())

	for confirmOption("Would you like to override another DNS record?", "") {
		spec = append(spec, getHostRecordPair())
	}

	return spec
}

func getDiskPressure() *v1beta1.DiskPressureSpec {
	if !confirmKind("Disk Pressure", "Applies IO pressure to the target") {
		return nil
	}

	spec := &v1beta1.DiskPressureSpec{}

	spec.Path = getInput(
		"Specify a path to apply IO pressure to",
		"Specify a specific mount point to target a specific disk",
	)

	if confirmOption("Would you like to apply read pressure?", "This applies read-based IO pressure (check the docs)") {
		readBPS, _ := strconv.Atoi(getInput("Specify the target amount of pressure, in bytes per second.", "check the docs"))
		spec.Throttling.ReadBytesPerSec = &readBPS
	}

	if confirmOption("Would you like to apply write pressure?", "This applies write-based IO pressure (check the docs)") {
		writeBPS, _ := strconv.Atoi(getInput("Specify the target amount of pressure, in bytes per second.", "check the docs"))
		spec.Throttling.WriteBytesPerSec = &writeBPS
	}

	return spec
}

func getCPUPressure() *v1beta1.CPUPressureSpec {
	if confirmKind("CPU Pressure", "Applies CPU pressure to the target") {
		return &v1beta1.CPUPressureSpec{}
	}

	return nil
}

func getNodeFailure() *v1beta1.NodeFailureSpec {
	if !confirmKind("Node Failure", "This will either shutdown or restart the targeted node (or node hosting the targeted pod)") {
		return nil
	}

	spec := &v1beta1.NodeFailureSpec{}
	spec.Shutdown = confirmOption("Would you like to shutdown the node permanently?",
		"Choosing yes will terminate the VM completely, and a new node will be scheduled. If you don't enable this, we will just restart the target node.")

	return spec
}

func getNetwork() *v1beta1.NetworkDisruptionSpec {
	if !confirmKind("Network Disruption", "Injects a variety of possible network issues") {
		return nil
	}

	spec := &v1beta1.NetworkDisruptionSpec{}

	// TODO rebase this onto main
	spec.Hosts = getSliceInput("Add a host to target (or leave blank)", "This will affect the network traffic between these hosts and your target.")
	spec.Port, _ = strconv.Atoi(getInput("what port to target (or leave blank)", "If specified, we will only affect traffic using this port"))

	spec.Flow, _ = selectInput(
		"Choose a flow direction",
		[]string{v1beta1.FlowEgress, v1beta1.FlowIngress},
		fmt.Sprintf("%s will affect traffic leaving the target. %s will not really affect traffic entering the target, but actually will affect replies to the inbound traffic.",
			v1beta1.FlowEgress, v1beta1.FlowIngress),
	)

	if confirmOption("Would you like to drop packets?", "Packets will be dropped before leaving the target") {
		spec.Drop, _ = strconv.Atoi(getInput("What % of packets should we affect?", "1-100%"))
	}

	if confirmOption("Would you like to duplicate packets?", "Packets will be duplicated immediately before leaving the target") {
		spec.Duplicate, _ = strconv.Atoi(getInput("What % of packets should we affect?", "1-100%"))
	}

	if confirmOption("Would you like to corrupt packets?", "Packets will be corrupted before leaving the target") {
		spec.Corrupt, _ = strconv.Atoi(getInput("What % of packets should we affect?", "1-100%"))
	}

	if confirmOption("Would you like to delay packets?", "Packets will be delayed before leaving the target") {
		delay, _ := strconv.ParseUint(getInput("How much to delay (in ms)?", "This will be the median amount of delay to apply"), 10, 0)
		spec.Delay = uint(delay)

		delayJitter, _ := strconv.ParseUint(getInput("What jitter on that delay (in ms)?", "This will be normally distributed around the delay you specified earlier. This will cause packets to re-order!"), 10, 0)
		spec.DelayJitter = uint(delayJitter)
	}

	if confirmOption("Would you like to limit bandwidth?", "bandwidthlimit") {
		spec.BandwidthLimit, _ = strconv.Atoi(getInput("What bandwidth limit should we set (in bytes per second)?", ">0"))
	}

	return spec
}

func getContainers() []string {
	if !confirmOption("Would you like to target a specific container[s]?", "The default is to target all containers in the target pod.") {
		return nil
	}

	containers := getSliceInput("Please enter a comma-delimited list of container name[s] to target.", "Please specify their names, not their IDs!")

	return containers
}

func getCount() *intstr.IntOrString {
	result := getInput(
		"How many targets would you like to disrupt? This can be an integer, or a percentage.",
		"Please specify an integer >0 or a percentage from 1% - 100%. If specifying a percentage, you must suffix with the % character, or we will think its an integer!",
	)
	// TODO somehow grab the other intrstr validate here

	wrappedResult := intstr.FromString(result)

	return &wrappedResult
}

func getSelectors() labels.Set {
	selectors := getSliceInput("Add a label selector[s] for targeting.", "Please specify this in the form of `key=value`, e.g., `app=hello-node`. One label selector per new-line. If you specify multiple, we will only target the union.")

	var selectorLabels labels.Set

	for _, s := range selectors {
		sAsSet, err := labels.ConvertSelectorToLabelsMap(s)

		if err != nil {
			fmt.Printf("invalid selector string: %v", err)
			return nil
		}

		selectorLabels = labels.Merge(selectorLabels, sAsSet)
	}

	return selectorLabels
}

func getLevel() types.DisruptionLevel {
	level, err := selectInput(
		"Select the Disruption Level.",
		[]string{types.DisruptionLevelNode, types.DisruptionLevelPod},
		"This will affect targeting with the label selectors, as well as injecting (depending on the disruption kind).",
	)

	if err != nil {
		level = types.DisruptionLevelPod
	}

	return types.DisruptionLevel(level)
}

func init() {}