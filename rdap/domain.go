package rdap

type Domain struct {
	ObjectClassName string       `json:"objectClassName"`
	Handle          string       `json:"handle"`
	LDHName         string       `json:"ldhName"`
	UnicodeName     string       `json:"unicodeName,omitempty"`
	Status          []string     `json:"status,omitempty"`
	Entities        []Entity     `json:"entities,omitempty"`
	Events          []Event      `json:"events,omitempty"`
	Nameservers     []Nameserver `json:"nameservers,omitempty"`
	SecureDNS       *SecureDNS   `json:"secureDNS,omitempty"`
}

type DomainResult struct {
	Domain Domain
	Raw    []byte
}

type Event struct {
	Action string `json:"eventAction"`
	Date   string `json:"eventDate"`
}

type Nameserver struct {
	ObjectClassName string `json:"objectClassName"`
	LDHName         string `json:"ldhName"`
	UnicodeName     string `json:"unicodeName,omitempty"`
}

type SecureDNS struct {
	DelegationSigned bool `json:"delegationSigned"`
}
