package supervisordreceiver

type Statistic struct {
	Name      string `xml:"name"`
	Start     string `xml:"start"`
	Stop      string `xml:"stop"`
	Now       string `xml:"now"`
	State     string `xml:"state"`
	StateName string `xml:"statename"`
	PID       int    `xml:"pid"`
}
