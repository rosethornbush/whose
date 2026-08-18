package rdap

type IPNetwork struct {
	ObjectClassName string   `json:"objectClassName"`
	Handle          string   `json:"handle,omitempty"`
	StartAddress    string   `json:"startAddress,omitempty"`
	EndAddress      string   `json:"endAddress,omitempty"`
	IPVersion       string   `json:"ipVersion,omitempty"`
	Name            string   `json:"name,omitempty"`
	Type            string   `json:"type,omitempty"`
	Country         string   `json:"country,omitempty"`
	Status          []string `json:"status,omitempty"`
	Entities        []Entity `json:"entities,omitempty"`
	Events          []Event  `json:"events,omitempty"`
}

type IPResult struct {
	Network IPNetwork
	Raw     []byte
}
