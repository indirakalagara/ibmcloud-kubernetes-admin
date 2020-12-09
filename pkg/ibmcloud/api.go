package ibmcloud

// TODO: return errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	// "os"
	// "io/ioutil"
)

// protocol
const protocol = "https://"

// subdomains
const (
	subdomainIAM                = "iam."
	subdomainAccounts           = "accounts."
	subdomainResourceController = "resource-controller."
	subdomainClusters           = "containers."
	subdomainUsers              = "users."
	subdomainTags               = "tags.global-search-tagging."
	subdomainBilling            = "billing."
)

// domain
const api = "cloud.ibm.com"


// endpoints
const (
	identityEndpoint     = protocol + subdomainIAM + api + "/identity/.well-known/openid-configuration"
	accountsEndpoint     = protocol + subdomainAccounts + api + "/coe/v2/accounts"
	resourcesEndpoint    = protocol + subdomainResourceController + api + "/v2/resource_instances"
	resourceKeysEndpoint = protocol + subdomainResourceController + api + "/v2/resource_keys"
	containersEndpoint   = protocol + subdomainClusters + api + "/global/v1"
	usersEndpoint        = protocol + subdomainUsers + api + "/v2"
	tagEndpoint          = protocol + subdomainTags + api + "/v3/tags"
	billingEndpoint      = protocol + subdomainBilling + api + "/v4/accounts"
	resourceEndoint      = protocol + subdomainResourceController + api + "/v1/resource_groups"
)

const (
	clusterEndpoint     = containersEndpoint + "/clusters"
	versionEndpount     = containersEndpoint + "/versions"
	locationEndpoint    = containersEndpoint + "/locations"
	zonesEndpoint       = containersEndpoint + "/zones"
	datacentersEndpoint = containersEndpoint + "/datacenters"
)

/*Manthan */
const (
	cfAppEndpoint     =  "/v3/apps/"
)

// grant types
const (
	passcodeGrantType     = "urn:ibm:params:oauth:grant-type:passcode"
	apikeyGrantType       = "urn:ibm:params:oauth:grant-type:apikey"
	refreshTokenGrantType = "refresh_token"
)

const basicAuth = "Basic Yng6Yng="

const cfToken="eyJhbGciOiJIUzI1NiIsImprdSI6Imh0dHBzOi8vdWFhLnVzLXNvdXRoLmNmLmNsb3VkLmlibS5jb20vdG9rZW5fa2V5cyIsImtpZCI6ImtleS0yIiwidHlwIjoiSldUIn0.eyJqdGkiOiJhZGVkMDNmMjNlYjg0OGVlYTRkYTg5MzQyN2FiM2Y4ZSIsInN1YiI6IjljYzg4OTVkLTAxYWItNGI3Yy1iMzRiLWYxZjBhM2ZiNTAwYyIsInNjb3BlIjpbIm9wZW5pZCIsIm5ldHdvcmsud3JpdGUiLCJ1YWEudXNlciIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm5ldHdvcmsuYWRtaW4iXSwiY2xpZW50X2lkIjoiY2YiLCJjaWQiOiJjZiIsImF6cCI6ImNmIiwiZ3JhbnRfdHlwZSI6InBhc3N3b3JkIiwidXNlcl9pZCI6IjljYzg4OTVkLTAxYWItNGI3Yy1iMzRiLWYxZjBhM2ZiNTAwYyIsIm9yaWdpbiI6IklCTWlkIiwidXNlcl9uYW1lIjoiaW5kaXJhLmthbGFnYXJhQGluLmlibS5jb20iLCJlbWFpbCI6ImluZGlyYS5rYWxhZ2FyYUBpbi5pYm0uY29tIiwiYXV0aF90aW1lIjoxNjA3MzMwMDM2LCJpYXQiOjE2MDczMzAwMzYsImV4cCI6MTYwNzMzMzYzNiwiaXNzIjoiaHR0cHM6Ly91YWEubmcuYmx1ZW1peC5uZXQvb2F1dGgvdG9rZW4iLCJ6aWQiOiJ1YWEiLCJhdWQiOlsiY2xvdWRfY29udHJvbGxlciIsInBhc3N3b3JkIiwiY2YiLCJ1YWEiLCJvcGVuaWQiLCJuZXR3b3JrIl0sImlhbV9pZCI6IklCTWlkLTA2MDAwMEg1NEMifQ.RbiT496dkjbY4nwAA_14O9h7qA6RLw-e04il1Q5H6Sw"

//// useful for loagging
// bodyBytes, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	panic(err)
// }
// bodyString := string(bodyBytes)
// log.Println(bodyString)
////

func timeTaken(t time.Time, name string) {
	elapsed := time.Since(t)
	log.Printf("TIME: %s took %s\n", name, elapsed)
}

func getError(resp *http.Response) error {
	var errorTemplate ErrorMessage
	if err := json.NewDecoder(resp.Body).Decode(&errorTemplate); err != nil {
		return err
	}
	if errorTemplate.Error != nil {
		return errors.New(errorTemplate.Error[0].Message)
	}
	if errorTemplate.Errors != nil {
		return errors.New(errorTemplate.Errors[0].Message)
	}
	return errors.New("unknown")
}

func getIdentityEndpoints() (*IdentityEndpoints, error) {
	result := &IdentityEndpoints{}
	err := fetch(identityEndpoint, nil, nil, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getToken(endpoint string, otp string) (*Token, error) {
	header := map[string]string{
		"Authorization": basicAuth,
	}

	form := url.Values{}
	form.Add("grant_type", passcodeGrantType)
	form.Add("passcode", otp)

	result := Token{}
	err := postForm(endpoint, header, nil, form, &result)

	if err != nil {
		log.Println("error in post form")
		return nil, err
	}

	return &result, nil
}

func getTokenFromIAM(endpoint string, apikey string) (*Token, error) {
	header := map[string]string{
		"Authorization": basicAuth,
	}

	form := url.Values{}
	form.Add("grant_type", apikeyGrantType)
	form.Add("apikey", apikey)

	result := &Token{}
	err := postForm(endpoint, header, nil, form, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func upgradeToken(endpoint string, refreshToken string, accountID string) (*Token, error) {
	header := map[string]string{
		"Authorization": basicAuth,
	}

	form := url.Values{}
	form.Add("grant_type", refreshTokenGrantType)
	form.Add("refresh_token", refreshToken)
	if accountID != "" {
		form.Add("bss_account", accountID)
	}

	result := &Token{}
	err := postForm(endpoint, header, nil, form, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getAccounts(endpoint *string, token string) (*Accounts, error) {
	if endpoint == nil {
		endpointString := accountsEndpoint
		endpoint = &endpointString
	} else {
		endpointString := accountsEndpoint + *endpoint
		endpoint = &endpointString
	}

	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	var result Accounts
	err := fetch(*endpoint, header, nil, &result)
	log.Println("In getAccounts : result ",result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getZones(showFlavors, location string) ([]Zone, error) {
	var result []Zone
	query := map[string]string{
		"showFlavors": showFlavors,
	}
	if len(location) > 0 {
		query["location"] = location
	}
	err := fetch(containersEndpoint+"/zones", nil, query, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getAccountResources(token, accountID string) (*AccountResources, error) {
	var result AccountResources
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	query := map[string]string{
		"account_id": accountID,
	}

	err := fetch(resourceEndoint, header, query, &result)
	if err != nil {
		return nil, err
	}
	//"/v1/resource_groups?account_id=9b13b857a32341b7167255de717172f5"
	return &result, nil
}

func getDatacenterVlan(token, refreshToken, datacenter string) ([]Vlan, error) {
	var result []Vlan
	header := map[string]string{
		"Authorization":        "Bearer " + token,
		"X-Auth-Refresh-Token": refreshToken,
	}

	url := datacentersEndpoint + "/" + datacenter + "/vlans"

	err := fetch(url, header, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getVersions() (*ClusterVersion, error) {
	var result ClusterVersion
	err := fetch(versionEndpount, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func getLocations() ([]Location, error) {
	var result []Location
	err := fetch(locationEndpoint, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getGeoLocations(geo string) ([]Location, error) {
	locations, err := getLocations()
	if err != nil {
		return nil, err
	}

	geoLocations := make([]Location, 0, 10)

	for _, location := range locations {
		if location.Kind == "dc" && location.Geography == geo {
			geoLocations = append(geoLocations, location)
		}
	}
	return geoLocations, nil
}

func getMachineTypes(datacenter, serverType, os string, cpuLimit, memoryLimit int) ([]MachineFlavor, error) {
	var result []MachineFlavor
	machineTypeEndpoint := fmt.Sprintf("%s/%s/machine-types", datacentersEndpoint, datacenter)
	err := fetch(machineTypeEndpoint, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	if serverType != "" && os != "" {
		filtered := make([]MachineFlavor, 0)
		toLower := strings.ToLower
		atoi := strconv.Atoi
		for _, machine := range result {
			cpu, _ := atoi(machine.Cores)
			memory, _ := atoi(strings.ReplaceAll(machine.Memory, "GB", ""))
			if toLower(machine.ServerType) == toLower(serverType) &&
				toLower(machine.Os) == toLower(os) &&
				cpu <= cpuLimit &&
				memory <= memoryLimit {
				filtered = append(filtered, machine)
			}
		}
		return filtered, nil
	}
	return result, nil
}

func getCluster(token, clusterID, resourceGroup string) (*Cluster, error) {
	var result Cluster
	header := map[string]string{
		"Authorization":         "Bearer " + token,
		"X-Auth-Resource-Group": resourceGroup,
	}
	err := fetch(clusterEndpoint+"/"+clusterID, header, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func getClusters(token, location string) ([]*Cluster, error) {
	defer timeTaken(time.Now(), "GetCluster :")
	var result []*Cluster
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	query := map[string]string{}
	if len(location) > 0 {
		query["location"] = location
	}

	err := fetch(clusterEndpoint, header, query, &result)

	if err != nil {
		return nil, err
	}

	// wg := &sync.WaitGroup{}

	// for _, cluster := range result {
	// 	time.Sleep(10 * time.Millisecond)
	// 	wg.Add(1)
	// 	go func(cluster *Cluster) {
	// 		tags, err := getTags(token, cluster.Crn)
	// 		if err != nil {
	// 			log.Println("error for tag: ", cluster.Name)
	// 			log.Println("error : ", err)
	// 		} else {
	// 			cluster.Tags = make([]string, len(tags.Items))
	// 			for i, val := range tags.Items {
	// 				cluster.Tags[i] = val.Name
	// 			}
	// 		}
	// 		wg.Done()
	// 	}(cluster)
	// 	wg.Add(1)
	// 	go func(cluster *Cluster) {
	// 		workers, err := getClusterWorkers(token, cluster.ID)
	// 		if err != nil {
	// 			log.Println("error for worker: ", cluster.Name)
	// 			log.Println("error : ", err)
	// 		} else {
	// 			cluster.Workers = workers
	// 			cost, err := getBillingData(token, accountID, cluster.Crn, workers)
	// 			if err != nil {
	// 				log.Println("error for cost: ", cluster.Name)
	// 			}
	// 			cluster.Cost = cost
	// 		}
	// 		wg.Done()
	// 	}(cluster)
	// }

	// wg.Wait()
	return result, nil
}

func getBillingData(token, accountID, clusterID, resourceInstanceID string) (string, error) {
	currentMonth := time.Now().Format("2006-01")
	workers, err := getClusterWorkers(token, clusterID)
	if err != nil {
		return "N/A", err
	}
	total := 0.0
	for _, worker := range workers {
		usage, err := getResourceUsagePerNode(token, accountID, currentMonth, resourceInstanceID, worker.ID)
		if err != nil {
			log.Printf("error getting resource usage %v\n", err)
			return "N/A", err
		}
		costForWorker := calcuateCostFromResourceUsage(usage)
		total += costForWorker
	}

	s := fmt.Sprintf("%.2f", total)

	return s, nil
}

func calcuateCostFromResourceUsage(usage *ResourceUsage) float64 {
	total := 0.0
	for _, resource := range usage.Resources {
		for _, use := range resource.Usage {
			total += use.Cost
		}
	}
	return total
}

func createCluster(token string, request CreateClusterRequest) (*CreateClusterResponse, error) {
	var result CreateClusterResponse
	header := map[string]string{
		"Authorization":         "Bearer " + token,
		"X-Auth-Resource-Group": request.ResourceGroup,
	}

	body, err := json.Marshal(request.ClusterRequest)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = postBody(clusterEndpoint, header, nil, body, &result)

	if err != nil {
		log.Println("error creating cluster : ", request.ClusterRequest.Name, err)
		return nil, err
	}
	log.Printf("cluster created. id :%s => name: %s", result.ID, request.ClusterRequest.Name)
	return &result, nil
}

func deleteCluster(token, id, resourceGroup, deleteResources string) error {
	header := map[string]string{
		"Authorization":         "Bearer " + token,
		"X-Auth-Resource-Group": resourceGroup,
	}

	query := map[string]string{
		"deleteResources": deleteResources,
	}

	deleteEndpoint := clusterEndpoint + "/" + id
	err := delete(deleteEndpoint, header, query, nil)
	if err != nil {
		return err
	}
	log.Println("cluster deleted, id :", id)
	return nil
}

func getClusterWorkers(token, id string) ([]Worker, error) {
	var result []Worker
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	workerEndpoint := clusterEndpoint + "/" + id + "/workers"

	err := fetch(workerEndpoint, header, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getResourceUsagePerNode(token, accountID, billingMonth, resourceInstanceID, workerID string) (*ResourceUsage, error) {
	var result ResourceUsage
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	crn := strings.ReplaceAll(resourceInstanceID, "::", ":worker:") + workerID
	query := map[string]string{
		"resource_id":          "containers-kubernetes",
		"_names":               "true",
		"resource_instance_id": crn,
	}

	endpoint := billingEndpoint + "/" + accountID + "/resource_instances/usage/" + billingMonth

	err := fetch(endpoint, header, query, &result)

	if err != nil {
		return nil, fmt.Errorf("error fetching resources usage %v", err)
	}

	return &result, err
}

func getTags(token string, crn string) (*Tags, error) {
	var result Tags
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	query := map[string]string{
		"attached_to": crn,
	}
	err := fetch(tagEndpoint, header, query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func setClusterTags(token, tag, clusterID, resourceGroup string) (*TagResult, error) {
	cluster, err := getCluster(token, clusterID, resourceGroup)
	if err != nil {
		log.Println("get cluster : ", err)
		return nil, err
	}
	crn := cluster.Crn

	resources := make([]Resource, 1)
	resources[0] = Resource{ResourceID: crn}
	updateTag := UpdateTag{TagName: tag, Resources: resources}
	tagResult, err := setTags(token, updateTag)
	if err != nil {
		log.Println("set tag : ", err)
		return nil, err
	}
	return tagResult, nil
}

func setTags(token string, updateTag UpdateTag) (*TagResult, error) {
	setTagsEndpoint := tagEndpoint + "/" + "attach"
	return updateTags(setTagsEndpoint, token, updateTag)
}

func deleteTags(token string, updateTag UpdateTag) (*TagResult, error) {
	setTagsEndpoint := tagEndpoint + "/" + "detach"

	return updateTags(setTagsEndpoint, token, updateTag)
}

func updateTags(endpoint, token string, updateTag UpdateTag) (*TagResult, error) {
	var result TagResult
	header := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	query := map[string]string{
		"providers": "ghost",
	}

	body, err := json.Marshal(updateTag)
	if err != nil {
		return nil, err
	}

	err = postBody(endpoint, header, query, body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

/* Manthan */
  
func initApplicationData(customerJSON string )([]byte, error) {
	return nil,nil
}

func getApplications(token, location string) ([]Application, error) {
	defer timeTaken(time.Now(), "GetApplication :")
	fmt.Println("In getApplications : Start ");
	//var returnResult []Application
	var result Applications
	// var app_resources []Application
	//app_result := {};
	
	/** Hard coded portion 
	appsData :=`{"resources":[{"guid":"cb90db45-a04f-405f-8180-3e916fab9a91","created_at":"2020-12-07T09:37:55Z","updated_at":"2020-08-16T15:44:16Z","name":"Manthan_VisualReg","state":"STARTED","lifecycle":{"type":"buildpack","data":{"buildpacks":[],"stack":"cflinuxfs3"}},"relationships":{"space":{"data":{"guid":"bd0c1abd-4d70-4e57-81c4-32af0fc43a00"}}},"metadata":{"labels":{},"annotations":{}}},
	{"guid":"cb90db45-a04f-405f-8180-3e916fab9a92","created_at":"2020-12-07T09:37:55Z","updated_at":"2020-08-16T15:44:16Z","name":"app2","state":"STOPPED","lifecycle":{"type":"buildpack","data":{"buildpacks":[],"stack":"cflinuxfs3"}},"relationships":{"space":{"data":{"guid":"bd0c1abd-4d70-4e57-81c4-32af0fc43a00"}}},"metadata":{"labels":{},"annotations":{}}},
	{"guid":"cb90db45-a04f-405f-8180-3e916fab9a93","created_at":"2020-12-07T09:37:55Z","updated_at":"2020-08-16T15:44:16Z","name":"app3","state":"STARTED","lifecycle":{"type":"buildpack","data":{"buildpacks":[],"stack":"cflinuxfs3"}},"relationships":{"space":{"data":{"guid":"bd0c1abd-4d70-4e57-81c4-32af0fc43a00"}}},"metadata":{"labels":{},"annotations":{}}},
	{"guid":"cb90db45-a04f-405f-8180-3e916fab9a94","created_at":"2020-12-07T09:37:55Z","updated_at":"2020-08-16T15:44:16Z","name":"app4","state":"STOPPED","lifecycle":{"type":"buildpack","data":{"buildpacks":[],"stack":"cflinuxfs3"}},"relationships":{"space":{"data":{"guid":"bd0c1abd-4d70-4e57-81c4-32af0fc43a00"}}},"metadata":{"labels":{},"annotations":{}}}
	]}`

	log.Println("In getApplications App Resources Results ", appsData);
	err := json.Unmarshal([]byte(appsData), &result)
	if err != nil {
		fmt.Println("Failed to unmarshal  result ")
	}
	**/
	
	// data, err := ioutil.ReadFile("/Users/indirakalagara/Workspaces/EB_Project/ibmcloud-kubernetes-admin/samples/cf_app_1.json")
	/*
	data, err := ioutil.ReadFile("cf_app_1.json")
    if err != nil {
      fmt.Print(err)
	} else {
		fmt.Print(data)
	}
	
	err = json.Unmarshal(data, &result)
    if err != nil {
        fmt.Println("error:", err)
	}
	*/

	token = cfToken;
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	query := map[string]string{}
	if len(location) > 0 {
		query["location"] = location
	}
	appEndPoint := "https://156dd4be-692f-48ec-afb3-3e0fb9464f0b.mock.pstmn.io/v1/applications"
	//"https://api.us-south.cf.cloud.ibm.com/v3/apps/";
	// err := fetch(appEndPoint, header, query, &(result.app_resources))
	err := fetch(appEndPoint, header, query, &result)
	if err != nil {
		log.Println("Error: In getApplications ", err);
		return nil, err
	}

	log.Println("In getApplications App Resources Results ", result);
	
	return result.Resources, nil
}

func getAppServiceBindings(token, app_guid string) ([]AppService, error) {
	defer timeTaken(time.Now(), "getAppServiceBindings :")
	fmt.Println("In getAppServiceBindings : Start ");
	//var returnResult []Application
	var result AppServices
	// var app_resources []Application
	/*
	appsData :=`{"resources":[{"guid":"2b762b65-020b-4afa-95a9-d4d88f20ffd3","type":"app","data":{"name":"manthan-visualreg1-visualrecogniti-1602496873610-5","instance_name":"manthan-visualreg1-visualrecogniti-1602496873610-5","binding_name":null,"credentials":{"redacted_message":"[PRIVATE DATA HIDDEN IN LISTS]"},"syslog_drain_url":null,"volume_mounts":[]},"created_at":"2020-10-12T10:14:02Z","updated_at":"2020-10-12T10:14:02Z","links":{"self":{"href":"https:\/\/api.us-south.cf.cloud.ibm.com\/v3\/service_bindings\/2b762b65-020b-4afa-95a9-d4d88f20ffd3"},"service_instance":{"href":"https:\/\/api.us-south.cf.cloud.ibm.com\/v2\/service_instances\/045a1ddb-89ba-4377-a0af-70b76257270d"},"app":{"href":"https:\/\/api.us-south.cf.cloud.ibm.com\/v3\/apps\/cb90db45-a04f-405f-8180-3e916fab9a91"}}},{"guid":"c6166b47-0f8e-4722-b079-1d5a11f9d9a7","type":"app","data":{"name":"Manthan-Db2-av","instance_name":"Manthan-Db2-av","binding_name":null,"credentials":{"redacted_message":"[PRIVATE DATA HIDDEN IN LISTS]"},"syslog_drain_url":null,"volume_mounts":[]},"created_at":"2020-10-20T08:43:21Z","updated_at":"2020-10-20T08:43:21Z","links":{"self":{"href":"https:\/\/api.us-south.cf.cloud.ibm.com\/v3\/service_bindings\/c6166b47-0f8e-4722-b079-1d5a11f9d9a7"},"service_instance":{"href":"https:\/\/api.us-south.cf.cloud.ibm.com\/v2\/service_instances\/e66ad0e9-7e4f-42eb-a78d-01c307169b14"},"app":{"href":"https:\/\/api.us-south.cf.cloud.ibm.com\/v3\/apps\/cb90db45-a04f-405f-8180-3e916fab9a91"}}}]}`

	log.Println("In getAppServiceBindings App Resources Results ", appsData);
	err := json.Unmarshal([]byte(appsData), &result)
	if err != nil {
		fmt.Println("Failed to unmarshal  result ")
	}
   */
   
 //  token = cfToken;
   header := map[string]string{
	   "Authorization": "Bearer " + token,
   }

	query := map[string]string{}
	log.Println("In getAppServiceBindings  app_guid  ", app_guid);
	if len(app_guid) > 0 {
		//log.Println("In getAppServiceBindings  app_guid  ", app_guid);
		query["app_guids"] = app_guid
	} else{
		query["app_guids"] = "cb90db45-a04f-405f-8180-3e916fab9a91"
	}

	appEndPoint := "https://156dd4be-692f-48ec-afb3-3e0fb9464f0b.mock.pstmn.io/v1/appserviceBindings";
	//"https://api.us-south.cf.cloud.ibm.com/v3/service_bindings";
	// err := fetch(appEndPoint, header, query, &(result.app_resources))
	err := fetch(appEndPoint, header, query, &result)
	if err != nil {
		log.Println("Error: In getAppServiceBindings ", err);
		return nil, err
	}
	
	// log.Println("In getAppServiceBindings  Results ", result);
	// log.Println("In getAppServiceBindings  Results Resources ", result.Resources);
	
	var appServiceslist []AppService;
	
	for _,tmpService := range result.Resources {
		tmpService.AppServiceName = tmpService.Data.Name
		tmpService.AppServiceInstanceName = tmpService.Data.InstanceName	
		// tmpService.Type ="service"
		appServiceslist = append(appServiceslist, tmpService)
		log.Println("In getAppServiceBindings  For  ", tmpService);
	
	}

	log.Println("In getAppServiceBindings  appServiceslist  ", appServiceslist);
	
	//return result.Resources, nil
	return appServiceslist, nil
}

// file, err := os.OpenFile("cf_apps_Manthan_servicesBindings.json", os.O_RDONLY, 0666)
	// 	//checkError(err)
	// 	b, err := ioutil.ReadAll(file)
	
	// if err != nil {
	// 	fmt.Print(err)
	// 	} else {
	// 		fmt.Print(b)
	// 	}
		
	// 	err = json.Unmarshal(b, &result)
	// 	if err != nil {
	// 		fmt.Println("error:", err)
	// 	}
	// 	log.Println("In getAppServiceBindings  b  ",  b);	
	// token = cfToken

	// header := map[string]string{
	// 	"Authorization": "Bearer " + token,
	// }