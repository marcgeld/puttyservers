package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"net/url"
	"os"
	"text/tabwriter"
)

const (
	HostnameKey = "HostName"
	registryKey = `Software\SimonTatham\PuTTY\Sessions`
)

func registryOpenKey(registryKey string) registry.Key {
	var access uint32 = registry.QUERY_VALUE | registry.ENUMERATE_SUB_KEYS
	regKey, err := registry.OpenKey(registry.CURRENT_USER, registryKey, access)
	if err != nil {
		if err == registry.ErrNotExist {
			log.Fatalf("Registry key '%s' not found %v\n", registryKey, err)
		}
		log.Fatalf("Failed to open '%s' Error: %v\n", registryKey, err)
	}
	return regKey
}

// Print (...to screen) or Write (...to file) json output
func printWriteJson(
	filename string,
	values map[string]string,
) {
	data, err := json.MarshalIndent(values, "", "\t")
	if err != nil {
		log.Fatal("could not create json", err)
	}
	if len(filename) > 0 {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			log.Fatalf("could not open file: %s %v", filename, err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatalf("could not close file: %s %v", filename, err)
			}
		}(f)
		_, err = f.Write(data)
		if err != nil {
			log.Fatal("could not write data to file ", err)
		}
		log.Printf("file '%s' written", filename)
	} else {
		fmt.Println(string(data))
	}
}

// Access Windows Registry for PuTTY Sessions
// Similar to: regedit /e "%USERPROFILE%\Desktop\putty-sessions.reg" HKEY_CURRENT_USER\Software\SimonTatham\PuTTY\Sessions
func main() {
	filename := flag.String("filename", "", "filename")
	isJson := flag.Bool("json", false, "Write json export")
	flag.Parse()
	log.Printf("Export Putty Session from registry")

	regKey := registryOpenKey(registryKey)
	defer func() {
		if err := regKey.Close(); err != nil {
			log.Fatalf("failed to close reg key '%v' error: %v", regKey, err)
		}
	}()

	keyNames, err := regKey.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatalf("Failed to get %q keys from registry error: %v", regKey, err)
	}

	var values = make(map[string]string)
	for _, encodedKey := range keyNames {
		subKeyString := fmt.Sprintf("%s\\%s", registryKey,
			encodedKey)
		subRegKey := registryOpenKey(subKeyString)
		hostname, _, err := subRegKey.GetStringValue(HostnameKey)
		if err != nil {
			log.Fatalf("Failed to get value for key '%s' from registry, error: %v", HostnameKey, err)
		}
		if err := subRegKey.Close(); err != nil {
			log.Fatalf("Failed to close reg key '%v' Error: %v\n", regKey, err)
		}
		// Don't add anything without hostname value.
		if len(hostname) > 1 {
			keyDecoded, _ := url.QueryUnescape(encodedKey)
			values[keyDecoded] = hostname
		}
	}

	if *isJson {
		printWriteJson(*filename, values)
	} else {
		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
		for key := range values {
			fmt.Fprintln(writer, fmt.Sprintf("%s\t%s", key, values[key]))
		}
		writer.Flush()
	}
}
