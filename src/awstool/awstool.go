/*
 * Assists in finding duplicates of Security Groups
 * Shows Security Groups
 * their network configs
 * and associated Instances
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Primary struct
type SecGroup struct {
	Id            string
	SecurityGroup ec2.SecurityGroup
	Instances     []ec2.Instance
}

// Sort implementation for by ec2 instances
type ByInstanceCount []SecGroup

func (g ByInstanceCount) Len() int {
	return len(g)
}
func (g ByInstanceCount) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}
func (g ByInstanceCount) Less(i, j int) bool {
	return len(g[i].Instances) > len(g[j].Instances)
}

// Sort implementation for by secgroup rules
type ByIPPort []ec2.IPPermission

func (p ByIPPort) Len() int {
	return len(p)
}
func (p ByIPPort) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p ByIPPort) Less(i, j int) bool {
	return *p[i].FromPort < *p[j].FromPort
}

var region string

func init() {
	regionFlag := flag.String("region", "us-east-1", "AWS Region")
	flag.Parse()
	region = string(*regionFlag)
}

func main() {

	log.Printf("AWS Region: %s", region)
	aws.DefaultConfig.Region = region

	// instances
	// security-groups
	//    with-delete

	deletestatements := false

	if contains(flag.Args(), "instances") {
		log.Println("Obtaining instances")
		// Describe Instances
		instancesResult, err := getInstances()
		if err != nil {
			log.Println("Can't even", err)
			return
		}
		log.Println("Obtained instances", len(instancesResult.Reservations))
		outputInstances(instancesResult)

	} else if contains(flag.Args(), "security-groups") {
		log.Print("Obtaining security groups")
		if contains(flag.Args(), "with-delete") {
			log.Println(" with delete statements")
			deletestatements = true
		}

		groups, err := getSecurityGroupsWithInstances()
		if err != nil {
			log.Println("Can't even", err)
			return
		}
		outputSecurityGroups(groups)

		if deletestatements {
			outputSecurityGroupsDeleteStatements(groups)
		}

	}

}

func getSecurityGroupsWithInstances() ([]SecGroup, error) {

	var groups []SecGroup

	grpmap := make(map[string]SecGroup)

	instancesResult, err := getInstances()
	if err != nil {
		return groups, err
	}

	svc := ec2.New(nil)

	// Security Groups
	params := &ec2.DescribeSecurityGroupsInput{}

	securityGroupsResult, err := svc.DescribeSecurityGroups(params)
	if err != nil {
		log.Println("Can't even", err)
		return groups, err
	}

	log.Println("Obtained security groups")

	for _, s := range securityGroupsResult.SecurityGroups {
		secgrp := SecGroup{SecurityGroup: *s, Id: *s.GroupID}
		id := *s.GroupID
		grpmap[id] = secgrp
	}

	// Associate Instances with Security Groups
	for _, r := range instancesResult.Reservations {
		for _, i := range r.Instances {
			for _, s := range i.SecurityGroups {
				secgrp := grpmap[*s.GroupID]
				secgrp.Instances = append(secgrp.Instances, *i)
				grpmap[*s.GroupID] = secgrp
			}

		}
	}

	// place array into named array struct for sorting purposes
	//var groups []SecGroup
	for _, e := range grpmap {
		groups = append(groups, e)
	}
	sort.Sort(ByInstanceCount(groups))

	return groups, nil

}

func outputSecurityGroups(groups []SecGroup) {
	// Output Security Groups
	fmt.Printf("%12s %20s %3s %3s %3s\n", "id", "name", "in", "out", "i")
	for _, v := range groups {
		fmt.Printf("%12s %20s %3v %3v %3v\n",
			*v.SecurityGroup.GroupID, *v.SecurityGroup.GroupName,
			len(v.SecurityGroup.IPPermissions), len(v.SecurityGroup.IPPermissionsEgress),
			len(v.Instances))

		if len(v.SecurityGroup.IPPermissions) > 0 {

			var ports []ec2.IPPermission
			for _, p := range v.SecurityGroup.IPPermissions {
				ports = append(ports, *p)
			}
			sort.Sort(ByIPPort(ports))

			for _, perm := range ports {
				if *perm.IPProtocol != "-1" {
					var cidrp string
					if len(perm.IPRanges) > 0 {
						cidrp = *perm.IPRanges[0].CIDRIP
					} else {
						cidrp = "all"
					}

					fmt.Printf("   %s %4v-%4v %s\n",
						*perm.IPProtocol, *perm.FromPort, *perm.ToPort,
						cidrp)
				}
			}
		}

		instances, _ := listInstances(v.Instances)
		if instances != "" {
			fmt.Printf("\tinstances: %s\n", instances)
		}
	}
}

func outputSecurityGroupsDeleteStatements(groups []SecGroup) {
	log.Println("AWS CLI to remove unused groups")
	fmt.Println()
	for _, d := range groups {
		if len(d.Instances) == 0 {
			fmt.Printf("aws ec2 delete-security-group --group-id %s --dry-run\n", *d.SecurityGroup.GroupID)
		}
	}
}

func getInstances() (ec2.DescribeInstancesOutput, error) {

	var instancesResult *ec2.DescribeInstancesOutput

	svc := ec2.New(nil)
	inParams := &ec2.DescribeInstancesInput{}
	instancesResult, err := svc.DescribeInstances(inParams)
	if err != nil {
		log.Println("Can't even", err)
		return *instancesResult, err
	}

	return *instancesResult, nil

}

func outputInstances(instancesResult ec2.DescribeInstancesOutput) {
	for _, r := range instancesResult.Reservations {
		fmt.Printf("Reservation %s, owner: %s\n", *r.ReservationID, *r.OwnerID)
		for _, i := range r.Instances {
			securityGroupsList, _ := listSecurityGroups(i.SecurityGroups)
			fmt.Printf("%s [%s]\n", *i.InstanceID, securityGroupsList)
		}
	}
}

func contains(array []string, item string) bool {
	set := make(map[string]struct{}, len(array))
	for _, s := range array {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// convenience method for concatenation, could probably do this more clever
func listInstances(instances []ec2.Instance) (string, error) {
	var iList string
	if len(instances) == 0 {
		return "", nil
	}
	if len(instances) == 1 {
		return *instances[0].InstanceID, nil
	}
	for _, v := range instances {
		iList += *v.InstanceID + ", "
	}
	iList = strings.TrimSuffix(iList, ", ")
	return iList, nil
}

// convenience method for concatenation, could probably do this more clever
func listSecurityGroups(groups []*ec2.GroupIdentifier) (string, error) {

	var groupList string

	if len(groups) == 1 {
		return *groups[0].GroupID, nil
	}

	for _, v := range groups {
		groupList += *v.GroupID + ", "
	}

	groupList = strings.TrimSuffix(groupList, ", ")

	return groupList, nil
}
