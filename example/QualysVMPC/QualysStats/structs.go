package main

import (
	"encoding/xml"
)

type ContactInfo struct {
	EMAIL string `xml:"EMAIL"`
}

type USER struct {
	XMLName       xml.Name    `xml:"USER"`
	Text          string      `xml:",chardata"`
	UserLogin     string      `xml:"USER_LOGIN"`
	UserId        string      `xml:"USER_ID"`
	CONTACT_INFO  ContactInfo `xml:"CONTACT_INFO"`
	UserStatus    string      `xml:"USER_STATUS"`
	UserLastLogin string      `xml:"LAST_LOGIN_DATE"`
	UserRole      string      `xml:"USER_ROLE"`
}

type USERLIST struct {
	XMLName   xml.Name `xml:"USER_LIST_OUTPUT"`
	Text      string   `xml:",chardata"`
	USER_LIST struct {
		USER []USER `xml:"USER"`
	} `xml:"USER_LIST"`
}

type UserRecord struct {
	Record    []string
	RowData   map[string]string
	FirstName string
	LastName  string
}

type ScanListResponse struct {
	Response struct {
		ScanList struct {
			Scan []struct {
				Ref    string `xml:"REF"`
				Title  string `xml:"TITLE"`
				Status struct {
					State    string `xml:"STATE"`
					SubState string `xml:"SUB_STATE"`
				} `xml:"STATUS"`
				LaunchDate         string `xml:"LAUNCH_DATETIME"`
				EndDate            string `xml:"END_DATETIME"`
				Duration           string `xml:"DURATION"`
				ProcessingPriority string `xml:"PROCESSING_PRIORITY"`
				Target             string `xml:"TARGET"`
			} `xml:"SCAN"`
		} `xml:"SCAN_LIST"`
	} `xml:"RESPONSE"`
}

type ApplianceListResponse struct {
	Response struct {
		ApplianceList struct {
			Appliance []struct {
				ID             string `xml:"ID"`
				Name           string `xml:"NAME"`
				SoftVersion    string `xml:"SOFTWARE_VERSION"`
				Status         string `xml:"STATUS"`
				MLVersion      string `xml:"ML_VERSION"`
				VulnSigVersion string `xml:"VULNSIGS_VERSION"`
				LastUpdate     string `xml:"LAST_UPDATED_DATE"`
				Type           string `xml:"TYPE"`
				Model          string `xml:"MODEL_NUMBER"`
				Serial         string `xml:"SERIAL_NUMBER"`
			} `xml:"APPLIANCE"`
		} `xml:"APPLIANCE_LIST"`
	} `xml:"RESPONSE"`
}

type Preferences struct {
	StartFromId  string `json:"startFromId"`
	LimitResults string `json:"limitResults"`
}

type Criteria struct {
	Value    string `xml:",chardata"`
	Field    string `xml:"field,attr"`
	Operator string `xml:"operator,attr"`
}

type Filter struct {
	Criteria Criteria `xml:"Criteria"`
}

type ServiceRequest struct {
	Preferences Preferences `json:"preferences"`
	Filters     Filter      `json:"filters"`
}

type Request struct {
	ServiceRequest ServiceRequest `xml:"ServiceRequest"`
}

type ServiceResponse struct {
	XMLName                   xml.Name `xml:"ServiceResponse"`
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	ResponseCode              string   `xml:"responseCode"`
	Count                     string   `xml:"count"`
	HasMoreRecords            string   `xml:"hasMoreRecords"`
	LastId                    string   `xml:"lastId"`
	Data                      struct {
		Text      string      `xml:",chardata"`
		HostAsset []HostAsset `xml:"HostAsset"`
	} `xml:"data"`
}

type HostAsset struct {
	Text     string `xml:",chardata"`
	ID       string `xml:"id"`
	Name     string `xml:"name"`
	Created  string `xml:"created"`
	Modified string `xml:"modified"`
	Type     string `xml:"type"`
	Tags     struct {
		Text string `xml:",chardata"`
		List struct {
			Text      string `xml:",chardata"`
			TagSimple struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id"`
				Name string `xml:"name"`
			} `xml:"TagSimple"`
		} `xml:"list"`
	} `xml:"tags"`
	CriticalityScore string `xml:"criticalityScore"`
	QwebHostId       string `xml:"qwebHostId"`
	LastVulnScan     string `xml:"lastVulnScan"`
	LastSystemBoot   string `xml:"lastSystemBoot"`
	LastLoggedOnUser string `xml:"lastLoggedOnUser"`
	Fqdn             string `xml:"fqdn"`
	Os               string `xml:"os"`
	DnsHostName      string `xml:"dnsHostName"`
	AgentInfo        struct {
		Text            string `xml:",chardata"`
		AgentVersion    string `xml:"agentVersion"`
		AgentId         string `xml:"agentId"`
		Status          string `xml:"status"`
		LastCheckedIn   string `xml:"lastCheckedIn"`
		ConnectedFrom   string `xml:"connectedFrom"`
		ChirpStatus     string `xml:"chirpStatus"`
		Platform        string `xml:"platform"`
		ActivatedModule string `xml:"activatedModule"`
		ManifestVersion struct {
			Text string `xml:",chardata"`
			Vm   string `xml:"vm"`
		} `xml:"manifestVersion"`
		AgentConfiguration struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id"`
			Name string `xml:"name"`
		} `xml:"agentConfiguration"`
		ActivationKey struct {
			Text         string `xml:",chardata"`
			ActivationId string `xml:"activationId"`
			Title        string `xml:"title"`
		} `xml:"activationKey"`
	} `xml:"agentInfo"`
	NetworkGuid     string `xml:"networkGuid"`
	Address         string `xml:"address"`
	TrackingMethod  string `xml:"trackingMethod"`
	Manufacturer    string `xml:"manufacturer"`
	Model           string `xml:"model"`
	TotalMemory     string `xml:"totalMemory"`
	Timezone        string `xml:"timezone"`
	BiosDescription string `xml:"biosDescription"`
	OpenPort        struct {
		Text string `xml:",chardata"`
		List struct {
			Text              string `xml:",chardata"`
			HostAssetOpenPort struct {
				Text        string `xml:",chardata"`
				Port        string `xml:"port"`
				Protocol    string `xml:"protocol"`
				ServiceName string `xml:"serviceName"`
			} `xml:"HostAssetOpenPort"`
		} `xml:"list"`
	} `xml:"openPort"`
	Software struct {
		Text string `xml:",chardata"`
		List struct {
			Text              string `xml:",chardata"`
			HostAssetSoftware struct {
				Text    string `xml:",chardata"`
				Name    string `xml:"name"`
				Version string `xml:"version"`
			} `xml:"HostAssetSoftware"`
		} `xml:"list"`
	} `xml:"software"`
	Vuln struct {
		Text string `xml:",chardata"`
		List struct {
			Text          string `xml:",chardata"`
			HostAssetVuln struct {
				Text               string `xml:",chardata"`
				Qid                string `xml:"qid"`
				HostInstanceVulnId string `xml:"hostInstanceVulnId"`
				FirstFound         string `xml:"firstFound"`
				LastFound          string `xml:"lastFound"`
			} `xml:"HostAssetVuln"`
		} `xml:"list"`
	} `xml:"vuln"`
	Processor struct {
		Text string `xml:",chardata"`
		List struct {
			Text               string `xml:",chardata"`
			HostAssetProcessor struct {
				Text  string `xml:",chardata"`
				Name  string `xml:"name"`
				Speed string `xml:"speed"`
			} `xml:"HostAssetProcessor"`
		} `xml:"list"`
	} `xml:"processor"`
	Volume struct {
		Text string `xml:",chardata"`
		List struct {
			Text            string `xml:",chardata"`
			HostAssetVolume struct {
				Text string `xml:",chardata"`
				Name string `xml:"name"`
				Size string `xml:"size"`
				Free string `xml:"free"`
			} `xml:"HostAssetVolume"`
		} `xml:"list"`
	} `xml:"volume"`
	Account struct {
		Text string `xml:",chardata"`
		List struct {
			Text             string `xml:",chardata"`
			HostAssetAccount struct {
				Text     string `xml:",chardata"`
				Username string `xml:"username"`
			} `xml:"HostAssetAccount"`
		} `xml:"list"`
	} `xml:"account"`
	NetworkInterface struct {
		Text string `xml:",chardata"`
		List struct {
			Text               string `xml:",chardata"`
			HostAssetInterface struct {
				Text           string `xml:",chardata"`
				Hostname       string `xml:"hostname"`
				InterfaceName  string `xml:"interfaceName"`
				MacAddress     string `xml:"macAddress"`
				Address        string `xml:"address"`
				GatewayAddress string `xml:"gatewayAddress"`
			} `xml:"HostAssetInterface"`
		} `xml:"list"`
	} `xml:"networkInterface"`
	IsDockerHost string `xml:"isDockerHost"`
}
type AllBinaryInfo struct {
	Text            string `xml:",chardata"`
	AgentId         string `xml:"agentId"`
	AgentVersion    string `xml:"agentVersion"`
	ManifestVersion string `xml:"manifestVersion"`
	Platform        string `xml:"platform"`
	Architecture    string `xml:"architecture"`
}

type Platform struct {
	Name        string `xml:"name"`
	Version     string `xml:"version"`
	Extenstions string `xml:"Extension"`
}

type Platforms struct {
	Text     string     `xml:",chardata"`
	Platform []Platform `xml:"Platform"`
}

type AgentType struct {
	XMLName      xml.Name `xml:"ServiceResponse"`
	Text         string   `xml:",chardata"`
	ResponseCode string   `xml:"responseCode"`
	Count        string   `xml:"count"`
	Data         struct {
		Text      string      `xml:",chardata"`
		Platforms []Platforms `xml:"platforms"`
	} `xml:"data"`
}

type BinaryInfoResponse struct {
	XMLName      xml.Name `xml:"ServiceResponse"`
	Text         string   `xml:",chardata"`
	ResponseCode string   `xml:"responseCode"`
	Count        string   `xml:"count"`
	Data         struct {
		Text          string `xml:",chardata"`
		AllBinaryInfo struct {
			Text      string `xml:",chardata"`
			Platforms struct {
				Text     string `xml:",chardata"`
				Platform []struct {
					Text      string `xml:",chardata"`
					Name      string `xml:"name"`
					Version   string `xml:"version,omitempty"`
					Hash      string `xml:"hash"`
					Extension string `xml:"extension"`
				} `xml:"Platform"`
			} `xml:"platforms"`
		} `xml:"AllBinaryInfo"`
	} `xml:"data"`
}
