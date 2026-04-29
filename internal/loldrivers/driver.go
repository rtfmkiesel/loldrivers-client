package loldrivers

// Based on the the JSON spec from
// https://github.com/magicsword-io/LOLDrivers/blob/main/bin/spec/drivers.spec.json

type Driver struct {
	ID                     string                  `json:"Id"`
	Author                 string                  `json:"Author"`
	Created                string                  `json:"Created"`
	MitreID                string                  `json:"MitreID"`
	Category               string                  `json:"Category"`
	Verified               string                  `json:"Verified"`
	Commands               Command                 `json:"Commands"`
	Resources              []string                `json:"Resources"`
	Acknowledgement        *Acknowledgement        `json:"Acknowledgement,omitempty"`
	Detection              []Detection             `json:"Detection,omitempty"`
	KnownVulnerableSamples []KnownVulnerableSample `json:"KnownVulnerableSamples,omitempty"`
	Tags                   []string                `json:"Tags"`
}

type Command struct {
	Command         string `json:"Command"`
	Description     string `json:"Description"`
	Usecase         string `json:"Usecase"`
	Privileges      string `json:"Privileges"`
	OperatingSystem string `json:"OperatingSystem"`
}

type Acknowledgement struct {
	Person string `json:"Person"`
	Handle string `json:"Handle"`
}

type Detection struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type KnownVulnerableSample struct {
	Filename string `json:"Filename"`
	MD5      string `json:"MD5,omitempty"`
	SHA1     string `json:"SHA1,omitempty"`
	SHA256   string `json:"SHA256,omitempty"`
}
