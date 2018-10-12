package explorer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// ResourceInfo ...
type ResourceInfo struct {
	TotalResource string `json:"totalResource"` // "2.35 DATX",		//已质押资源
	CPUResource   string `json:"cpuResource"`   // "1.75 DATX",
	CPUUsed       string `json:"cpuUsed"`       // "2.25 ms",		//cpu资源已用
	CPUTotal      string `json:"cpuTotal"`      // "29.42 ms",		//cpu资源总量
	NetResource   string `json:"netResource"`   // "1.6 DATX",
	NetUsed       string `json:"netUsed"`       // "2.25 KB",		//net资源已用
	NetTotal      string `json:"netTotal"`      // "29.42 KB",		//net资源总量
	RAMResource   string `json:"ramResource"`   // "1.6 DATX",
	RAMUsed       string `json:"ramUsed"`       // "15 KB",			//ram资源已用
	RAMTotal      string `json:"ramTotal"`      // "718.2 KB"		//ram资源总量
}

// DATXResource ...
type DATXResource struct {
	NetWeight int64 `json:"net_weight"`
	CPUWeight int64 `json:"cpu_weight"`
	NetLimit  struct {
		Used int64  `json:"used"`
		Max  string `json:"max"`
	} `json:"net_limit"`
	CPULimit struct {
		Used int64  `json:"used"`
		Max  string `json:"max"`
	} `json:"cpu_limit"`
	RAMUsage       int64 `json:"ram_usage"`
	TotalResources struct {
		NetWeight string `json:"net_weight"`
		CPUWeight string `json:"cpu_weight"`
		RAMBytes  int64  `json:"ram_bytes"`
	} `json:"total_resources"`
}

// DATXGetResourceFormData ...
type DATXGetResourceFormData struct {
	Account string `json:"account_name"`
}

// DATXResourceRequest ...
type DATXResourceRequest struct {
	Account string `form:"account" json:"account" binding:"required"`
}

// GetDATXResource ...
func GetDATXResource(account string) (*ResourceInfo, error) {
	formData := DATXGetResourceFormData{account}
	bytesData, err := json.Marshal(formData)
	if err != nil {
		return nil, fmt.Errorf("datx get_account parameter error %v", formData)
	}

	URL := "http://172.31.3.38:8888/v1/chain/get_account"
	request, err := http.NewRequest("POST", URL, bytes.NewReader(bytesData))
	if err != nil {
		return nil, fmt.Errorf("datx get_account request error %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("datx get_account resp error %v", err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, fmt.Errorf("datx get_account not found %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("datx get_account body error %v", err)
	}

	resource := &DATXResource{}
	if err := json.Unmarshal([]byte(body), &resource); err != nil {
		return nil, fmt.Errorf("datx get_actions unmarshal error %v", string(body))
	}

	info := &ResourceInfo{}

	cpuUsed := float64(resource.CPULimit.Used)
	cpuTotal, _ := strconv.ParseFloat(resource.CPULimit.Max, 64)
	netUsed := float64(resource.NetLimit.Used)
	netTotal, _ := strconv.ParseFloat(resource.NetLimit.Max, 64)
	ramUsed := float64(resource.RAMUsage)
	ramTotal := float64(resource.TotalResources.RAMBytes)

	cpuResource := float64(resource.CPUWeight) / 10000
	netResource := float64(resource.NetWeight) / 10000
	ramResource := ramTotal / (64*1024*1024*1024 + 1) * 10000000000

	info.TotalResource = strconv.FormatFloat(cpuResource+netResource+ramResource, 'f', 4, 64) + " DATX"
	info.CPUResource = resource.TotalResources.CPUWeight
	info.CPUUsed = strconv.FormatFloat(cpuUsed/1000, 'f', 2, 64) + " ms"
	info.CPUTotal = strconv.FormatFloat(cpuTotal/1000, 'f', 2, 64) + " ms"
	info.NetResource = resource.TotalResources.NetWeight
	info.NetUsed = strconv.FormatFloat(netUsed/1024, 'f', 2, 64) + " KB"
	info.NetTotal = strconv.FormatFloat(netTotal/1024, 'f', 2, 64) + " KB"
	info.RAMResource = strconv.FormatFloat(ramResource, 'f', 4, 64) + " DATX"
	info.RAMUsed = strconv.FormatFloat(ramUsed/1024, 'f', 2, 64) + " KB"
	info.RAMTotal = strconv.FormatFloat(ramTotal/1024, 'f', 2, 64) + " KB"

	return info, nil
}
