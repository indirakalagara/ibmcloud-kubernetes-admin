package server

import (
	"encoding/json"
	// "fmt"
	"net/http"

	// "github.com/gorilla/mux"
	// "github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)


func (s *Server) AppListHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	applications, err := session.GetApplications("")
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get applications", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(applications)
}

func (s *Server) AppServiceBindingsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	app_guid := r.FormValue("app_guid")
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
	// funcVersion := r.FormValue("version")
	// if err != nil {
	// 	handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
	// 	return
	// }

	applications, err := session.GetAppServiceBindings(app_guid)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get applications services", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(applications)
}