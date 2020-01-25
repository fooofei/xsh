package base

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"regexp"
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

func initXHost() {
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
				if isUseRange(address) {
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

func isUseRange(name string) bool {
	if ok, _ := regexp.MatchString("^[0-9.]+-[0-9.]+$", name); !ok {
		return false
	}
	return true
}

//64.233.196.0-64.233.196.25
func processIpRange(ip string) ([]string, error) {
	fields := strings.Split(ip, "-")
	if !CheckIp(fields[0]) || !CheckIp(fields[1]) {
		log.Fatalf("ip[%s] illegal", ip)
	}

	start := ip2num(fields[0])
	end := ip2num(fields[1])

	if start >= end {
		log.Fatalf("ip[%s] illegal", ip)
	}

	var ret []string
	for i := start; i <= end; i++ {
		if i&0xff > 0 {
			ret = append(ret, num2ip(i))
		}
	}
	return ret, nil
}

func ip2num(ip string) int {
	fields := strings.Split(ip, ".")
	a, _ := strconv.Atoi(fields[0])
	b, _ := strconv.Atoi(fields[1])
	c, _ := strconv.Atoi(fields[2])
	d, _ := strconv.Atoi(fields[3])
	return a<<24 | b<<16 | c<<8 | d
}

func num2ip(num int) string {
	a := (num >> 24) & 0xff
	b := (num >> 16) & 0xff
	c := (num >> 8) & 0xff
	d := num & 0xff
	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}
