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

type Host struct {
	Name      string       `yaml:"name"`
	Auth      string       `yaml:"auth,omitempty"`
	Addresses []string     `yaml:"addresses,omitempty"`
	Groups    []string     `yaml:"groups,omitempty"`
	Details   []HostDetail `yaml:"details,omitempty"`
	Port      int          `yaml:"port,omitempty"`
	AllHost   []HostDetail
}

type xHost struct {
	Hosts []Host `yaml:"hosts,omitempty"`
}

var XHost = xHost{}
var XHostMap = make(map[string]Host)

func InitXHost() {
	var filePath = path.Join(RootPath, HostFile)

	h, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Can not read host file[%s].", filePath)
	}
	err = yaml.Unmarshal(h, &XHost)
	if err != nil {
		log.Fatalf("host file[%s] unmarshal error: %v", filePath, err)
	}

	if len(XHost.Hosts) == 0 {
		log.Fatal("The hosts empty.")
	}

	for _, value := range XHost.Hosts {
		if !CheckName(value.Name) {
			log.Fatalf("Host name [%s] illegal", value.Name)
		}
		XHostMap[value.Name] = value
	}

	if len(XHost.Hosts) != len(XHostMap) {
		log.Fatal("Host duplicate")
	}

	postProcessHostGroup()
}

func newHostDetail(address string, hostGroup Host) HostDetail {
	auth := XAuthMap[hostGroup.Auth]

	result := HostDetail{
		Address:    address,
		Port:       0,
		Username:   auth.Username,
		Password:   auth.Password,
		PrivateKey: auth.PrivateKey,
		Passphrase: auth.Passphrase,
		SuType:     auth.SuType,
		SuPass:     auth.SuPass,
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
	var err error
	for name, xHost := range XHostMap {
		allHost := make([]HostDetail, 0)
		for _, value := range xHost.Details {
			if !checkHostDetail(value) {
				log.Fatalf("Host[%s] detail[%s] illegal", name, value.Address)
			}
			if value.Password, err = GetPlainPassword(value.Password); err != nil {
				Error.Printf("host[%s] detail[%s] password decrypt error: %v", name, value.Address, err)
			}
			if value.Passphrase, err = GetPlainPassword(value.Passphrase); err != nil {
				Error.Printf("host[%s] detail[%s] passphrase decrypt error: %v", name, value.Address, err)
			}
			if value.SuPass, err = GetPlainPassword(value.SuPass); err != nil {
				Error.Printf("host[%s] detail[%s] suPass decrypt error: %v", name, value.Address, err)
			}
			allHost = append(allHost, value)
		}

		if len(xHost.Addresses) != 0 {
			for _, address := range xHost.Addresses {
				if strings.Contains(address, "-") {
					if ips, err := processIpRange(address); err != nil {
						log.Fatalf("Hostgroup[%s] Addresses[%s] illegal, err: %v", name, address, err)
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
		if len(xHost.Groups) != 0 {
			for _, hostGroup := range xHost.Groups {
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
			log.Fatalf("Hostgroup[%s] Addresses[%s] illegal", name, xHost.AllHost[0].Address)
		}

		for i := 1; i < len(xHost.AllHost); i++ {
			if xHost.AllHost[i] == xHost.AllHost[i-1] {
				log.Fatalf("Hostgroup[%s] Addresses[%s] duplicate", name, xHost.AllHost[i].Address)
			}
			if !checkHostDetail(xHost.AllHost[i]) {
				log.Fatalf("Hostgroup[%s] Addresses[%s] illegal", name, xHost.AllHost[i].Address)
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
