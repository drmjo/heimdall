package jsonUtil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

var hasRefs bool

/**
 * for the grep to work properly the json needs to be minified
 */
func UglifyJson(refJson []byte) ([]byte, error) {
	var rj interface{}
	err := json.Unmarshal(refJson, &rj)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	return json.Marshal(rj)
}

/**
 * save data to file api.json.resolved.json
 */
func WriteToFile(data []byte, rootJsonPath string, pretty bool) {

	if pretty {
		log.Printf("Prettyfying the json before write.")
		var tmp interface{}
		json.Unmarshal(data, &tmp)
		pretty, err := json.MarshalIndent(&tmp, "", " ")
		if err != nil {
			log.Printf("Error prettyfying json: %v", err.Error())
		}
		data = pretty
	}

	filename := fmt.Sprintf("%s.resolved.json", rootJsonPath)
	log.Printf("Writing to file: %s", filename)
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		log.Fatal(err.Error())
	}
}

func resolveRefs(json []byte, baseDir string) []byte {
	refExpression := regexp.MustCompile(`(?U){"\$ref":"(.*)"}`)

	if !refExpression.Match(json) {
		log.Print("Does not have refs returning untouched")
		return json
	}

	log.Print("Has Refs should be resolved")
	parsed := resolveRefs(refExpression.ReplaceAllFunc(json, func(match []byte) []byte {
		refUrl := refExpression.FindSubmatch(match)[1]
		url, err := url.Parse(string(refUrl))
		if err != nil {
			log.Print(err.Error())
			return match
		}
		refJsonPath := fmt.Sprintf("%s/%s", baseDir, url.EscapedPath())
		refJson, err := ioutil.ReadFile(refJsonPath)
		if err != nil {
			log.Print(err.Error())
			return match
		}
		refJson, _ = UglifyJson(refJson)
		return refJson
	}), baseDir)

	return parsed
}

func ResolveRefs(rootJsonPath string) []byte {
	rootJson, err := ioutil.ReadFile(rootJsonPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rootJson, err = UglifyJson(rootJson)
	if err != nil {
		log.Fatal(err.Error())
	}
	parsed := resolveRefs(rootJson, filepath.Dir(rootJsonPath))
	return parsed
}
