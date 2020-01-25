package base

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"sort"
	"strconv"
	"strings"
)

type HostDetail struct {
	Address    string `yaml:"address,omitempty"`
	Port       int    `yaml:"port,omitempty"`
	Username   string `yaml:"username,omitempty"`
	Password   string `yaml:"password,omitempty"`
	PrivateKey string `yaml:"privatekey,omitempty"`
	Passphrase string `yaml:"passphrase,omitempty"`
	SuType     string `yaml:"sutype,omitempty"`
	SuPass     string `yaml:"supass,omitempty"`
}

type HostGroup struct {
	Name           string       `yaml:"name,omitempty"`
	Authentication string       `yaml:"authentication,omitempty"`
	HostAddresses  []string     `yaml:"hostaddresses,omitempty"`
	HostGroups     []string     `yaml:"hostgroups,omitempty"`
	HostDetails    []HostDetail `yaml:"hostdetails,omitempty"`
	Port           int          `yaml:"port,omitempty"`
	AllHost        []HostDetail
}

type xHost struct {
	HostsGroups []HostGroup `yaml:"hostgroups,omitempty"`
}

var XHost = xHost{}
var XHostMap = make(map[string]HostGroup)

func InitXHost() {
	var filePath = path.Join(ConfigRootPath, HostgroupsFile)

	h, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Can not read hostgroups file[%s].", filePath)
	}
	err = yaml.Unmarshal(h, &XHost)
	if err != nil {
		log.Fatalf("Hostgroups[%s] unmarshal error: %v", filePath, err)
	}

	if len(XHost.HostsGroups) == 0 {
		log.Fatal("The hostgroups empty.")
	}

	for _, value := range XHost.HostsGroups {
		if !CheckName(value.Name) {
			log.Fatalf("Hostgroup name [%s] illegal", value.Name)
		}
		XHostMap[value.Name] = value
	}

	if len(XHost.HostsGroups) != len(XHostMap) {
		log.Fatal("Hostgroup duplicate")
	}

	postProcessHostGroup()
}

func newHostDetail(address string, hostGroup HostGroup) HostDetail {
	authentication := XAuthMap[hostGroup.Authentication]

	result := HostDetail{
		Address:    address,
		Port:       0,
		Username:   authentication.Username,
		Password:   authentication.Password,
		PrivateKey: authentication.PrivateKey,
		Passphrase: authentication.Passphrase,
		SuType:     authentication.SuType,
		SuPass:     authentication.SuPass,
	}
	if hostGroup.Port > 0 {
		result.Port = hostGroup.Port
	}

	return result
}

func checkHostDetail(hostDetail HostDetail) bool {
	return hostDetail.Address != ""
}

func postProcessHostGroup() {
	for name, xHost := range XHostMap {
		allHost := make([]HostDetail, 0)
		for _, v := range xHost.HostDetails {
			if !checkHostDetail(v) {
				log.Fatalf("Hostgroup[%s] HostDetails[%s] illegal", name, v.Address)
			}
			v.Password = GetPlainPassword(v.Password)
			v.SuPass = GetPlainPassword(v.SuPass)
			v.Passphrase = GetPlainPassword(v.Passphrase)
			allHost = append(allHost, v)
		}

		if len(xHost.HostAddresses) != 0 {
			for _, address := range xHost.HostAddresses {
				if strings.Contains(address, "-") {
					if ips, err := processIpRange(address); err != nil {
						log.Fatalf("Hostgroup[%s] HostAddresses[%s] illegal, err: %v", name, address, err)
					} else {
						for _, ip := range ips {
							allHost = append(allHost, newHostDetail(ip, xHost))
						}
					}
				} else {
					allHost = append(allHost, newHostDetail(address, xHost))
				}
			}
		}

		if len(allHost) == 0 {
			log.Fatalf("Hostgroup[%s] empty", name)
		}
		xHost.AllHost = allHost

		XHostMap[name] = xHost
	}

	for name, xHost := range XHostMap {
		if len(xHost.HostGroups) != 0 {
			for _, hostGroup := range xHost.HostGroups {
				h, ok := XHostMap[hostGroup]
				if !ok {
					log.Fatalf("Hostgroup[%s] -> Hostgroup[%s] not found", name, hostGroup)
				}
				xHost.AllHost = append(xHost.AllHost, h.AllHost...)
			}
		}

		sort.Slice(xHost.AllHost, func(i, j int) bool {
			return xHost.AllHost[i].Address < xHost.AllHost[j].Address
		})
		if !checkHostDetail(xHost.AllHost[0]) {
			log.Fatalf("Hostgroup[%s] HostAddresses[%s] illegal", name, xHost.AllHost[0].Address)
		}

		for i := 1; i < len(xHost.AllHost); i++ {
			if xHost.AllHost[i] == xHost.AllHost[i-1] {
				log.Fatalf("Hostgroup[%s] HostAddresses[%s] duplicate", name, xHost.AllHost[i].Address)
			}
			if !checkHostDetail(xHost.AllHost[i]) {
				log.Fatalf("Hostgroup[%s] HostAddresses[%s] illegal", name, xHost.AllHost[i].Address)
			}
		}

		XHostMap[name] = xHost
	}
}

//64.233.196.0-64.233.196.25
func processIpRange(ip string) ([]string, error) {
	fields := strings.Split(ip, "-")
	if !CheckIp4(fields[0]) || !CheckIp4(fields[1]) {
		log.Fatalf("ip range[%s] format illegal", ip)
	}

	sFields := strings.Split(fields[0], ".")
	eFields := strings.Split(fields[1], ".")

	s0 := sFields[0]
	s1 := sFields[1]
	s2 := sFields[2]
	s3, _ := strconv.Atoi(sFields[3])

	e0 := eFields[0]
	e1 := eFields[1]
	e2 := eFields[2]
	e3, _ := strconv.Atoi(eFields[3])

	if s0 != e0 || s1 != e1 || s2 > e2 {
		log.Fatalf("ip range[%s] illegal", ip)
	}

	if s2 == e2 && s3 >= e3 {
		log.Fatalf("ip range[%s] illegal", ip)
	}

	var ret []string

	if s2 == e2 {
		for i := s3; i <= e3; i++ {
			ret = append(ret, s0+"."+s1+"."+s2+"."+strconv.Itoa(i))
		}
	} else {
		s2, _ := strconv.Atoi(sFields[2])
		e2, _ := strconv.Atoi(eFields[2])

		for i := s3; i <= 255; i++ {
			ret = append(ret, s0+"."+s1+"."+strconv.Itoa(s2)+"."+strconv.Itoa(i))
		}
		for i := s2 + 1; i < e2; i++ {
			for j := 0; j <= 255; j++ {
				ret = append(ret, s0+"."+s1+"."+strconv.Itoa(i)+"."+strconv.Itoa(j))
			}
		}
		for i := 0; i <= e3; i++ {
			ret = append(ret, s0+"."+s1+"."+strconv.Itoa(e2)+"."+strconv.Itoa(i))
		}
	}

	return ret, nil
}
