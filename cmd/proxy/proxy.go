package main

import (
    "net/http"
    "fmt"
    "os"
    "net/url"
    "encoding/json"
    "xmljsonproxy"
)

func writeJson(w http.ResponseWriter, obj interface{}, code int) error {
	marshalled, err := json.MarshalIndent(map[string]interface{}{"result":obj}, "", "  ")
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("%s\n", marshalled)))

	return nil
}

var ccpEndpoint, _ = url.Parse("https://api.eveonline.com")

type server struct {
}


func makeRequest(w http.ResponseWriter, r *http.Request) {
    targetUrl := ccpEndpoint.ResolveReference(r.URL)
    resp, err := http.Get(targetUrl.String())
    
    if err != nil {
        fmt.Println(err)
        writeJson(w, err.Error(), 400)
        return
    }
    
    result, err := xmljsonproxy.Transform(resp.Body)
    
    if err != nil {
        fmt.Println(err)
        writeJson(w, err.Error(), 400)
        return
    }
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(resp.StatusCode)
    w.Write(result)
    
    return    
}


func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {    
    makeRequest(w, r)
    
    return
}

func main() {
    r := server{}
	if port := os.Getenv("PORT"); port != "" {
	    http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	} else {
	    http.ListenAndServe(":9293", r)
	}
}