package event

import (
	"common"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"os/exec"
	"strings"
)

//监听到容器启动时做的操作
func Start(client *docker.Client, event *docker.APIEvents) {
	fmt.Println("Received start event %s for container %s", event.Status, event.ID[:12])
	container, err := client.InspectContainer(event.ID[:12])
	common.ErrorHandle(err)
	ipAddress := container.NetworkSettings.IPAddress
	include := common.GetConfig("Section", "include")
	if include == "*" || strings.Contains(","+include+",", ","+common.SubstrAfter(container.Name, 0)+",") {
		common.AppendFile(common.GetConfig("Section", "hostFile"), ipAddress+"  "+getDomainName(container.Name))
		restartDns()
	}
}

//监听到容器消亡时做的操作
func Die(client *docker.Client, event *docker.APIEvents) {
	fmt.Println("Received die event %s for container %s", event.Status, event.ID[:12])
	container, err := client.InspectContainer(event.ID[:12])
	common.ErrorHandle(err)
	include := common.GetConfig("Section", "include")
	if include == "*" || strings.Contains(","+include+",", ","+common.SubstrAfter(container.Name, 0)+",") {
		strData := common.ReadFile(common.GetConfig("Section", "hostFile"))
		arrData := strings.Split(strData, "\n")
		strData = ""
		for i := 0; i < len(arrData); i++ {
			if strings.Index(arrData[i], getDomainName(container.Name)) >= 0 {
				continue
			}
			if strData == "" {
				strData = arrData[i]
			} else {
				strData += "\n" + arrData[i]
			}
		}
		common.SaveFile(common.GetConfig("Section", "hostFile"), strData)
		restartDns()
	}
}

//重启dnsmasq
func restartDns() {
	sh := exec.Command("/sbin/service", "dnsmasq", " restart")
	_, err := sh.CombinedOutput()
	common.ErrorHandle(err)
}

//根据容器名称生产域名
func getDomainName(cName string) string {
	cName = common.SubstrAfter(cName, 0)
	index := strings.Index(cName, "_")
	if index < 0 {
		return cName + ".com"
	}
	lcName := common.SubstrBefore(cName, index)
	rcName := common.SubstrAfter(cName, index)
	index = strings.Index(rcName, "_")
	if index < 0 {
		return lcName + "." + rcName + ".com"
	}
	lrcName := common.SubstrBefore(rcName, index)
	return lcName + "." + lrcName + ".com"
}
