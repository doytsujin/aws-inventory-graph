package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type availabilityZoneList struct {
	*ec2.DescribeAvailabilityZonesOutput
}
type availabilityZoneNodes []availabilityZoneNode

type availabilityZoneNode struct {
	UID      string   `json:"uid,omitempty"`
	Type     []string `json:"dgraph.type,omitempty"`
	Name     string   `json:"name,omitempty"` // This field is only for Ratel Viz
	Region   string   `json:"Region,omitempty"`
	Service  string   `json:"Service,omitempty"`
	State    string   `json:"State,omitempty"`
	ZoneName string   `json:"ZoneName,omitempty"`
	ZoneID   string   `json:"ZoneId,omitempty"`
}

func (c *connector) listAvailabilityZones() availabilityZoneList {
	defer c.waitGroup.Done()

	log.Println("List AvailibiltyZones")
	response, err := ec2.New(c.awsSession).DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		log.Fatal(err)
	}
	return availabilityZoneList{response}
}

func (list availabilityZoneList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.AvailabilityZones) == 0 {
		return
	}
	log.Println("Add AvailibilityZones Nodes")
	a := make(availabilityZoneNodes, 0, len(list.AvailabilityZones))

	for _, i := range list.AvailabilityZones {
		var b availabilityZoneNode
		b.Service = "ec2"
		b.Type = []string{"AvailabilityZone"}
		b.Region = c.awsRegion
		b.Name = *i.ZoneName
		b.ZoneName = *i.ZoneName
		b.ZoneID = *i.ZoneId
		a = append(a, b)
	}
	c.dgraphAddNodes(a)

	m := make(map[string]availabilityZoneNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("AvailabilityZone"), &m)
	for _, i := range m["list"] {
		n[i.ZoneName] = i.UID
	}
	c.ressources["AvailabilityZones"] = n
}