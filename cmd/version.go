package cmd

var Version = "1.2.1"


type Releases struct {
	Releases  []Release       `json:"releases"`
}

type Release struct {
	Version string        `json:"version"`
}

func getVersion() string {

     


}
