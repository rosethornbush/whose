package rdap

type Autnum struct {
	ObjectClassName string   `json:"objectClassName"`
	Handle          string   `json:"handle,omitempty"`
	StartAutnum     uint32   `json:"startAutnum,omitempty"`
	EndAutnum       uint32   `json:"endAutnum,omitempty"`
	Name            string   `json:"name,omitempty"`
	Type            string   `json:"type,omitempty"`
	Country         string   `json:"country,omitempty"`
	Status          []string `json:"status,omitempty"`
	Entities        []Entity `json:"entities,omitempty"`
	Events          []Event  `json:"events,omitempty"`
}

type ASNResult struct {
	Autnum Autnum
	Raw    []byte
}
