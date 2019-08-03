// +build ignore

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func getURLStrings(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []string
	reader := bufio.NewReader(resp.Body)
	for {
		s, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err == io.EOF {
			break
		}

		s = strings.Trim(s, "\n\t ")
		if len(s) < 1 {
			continue
		}

		if s[0] != '#' {
			data = append(data, strings.ToLower(strings.Trim(s, "',")))
		}
	}

	return data, nil

}

func getURLJson(url string, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}

	return nil
}

func getURLPHP(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []string
	reader := bufio.NewReader(resp.Body)
	for {
		s, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err == io.EOF {
			break
		}

		s = strings.Trim(s, "\n\t ")
		if len(s) < 1 {
			continue
		}

		if s[0] == '\'' || s[0] == '"' {
			data = append(data, strings.Trim(s, `"',`))
		}
	}

	return data, nil
}

func printMap(w io.Writer, name string, data []string) {
	_, _ = fmt.Fprintf(w, "\nvar %s = map[string]bool {\n", name)
	for i := range data {
		_, _ = fmt.Fprintf(w, "\t%q: true,\n", data[i])
	}
	fmt.Fprintln(w, "}\n")
}

func formatCode(src []byte, out io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return err
	}

	return format.Node(out, fset, f)
}

func main() {
	f := flag.String("file", "data.go", "file to generate")
	flag.Parse()

	fl, err := os.Create(*f)
	if err != nil {
		log.Print(err)
		return
	}
	defer fl.Close()

	var disposable []string
	if err := getURLJson("https://raw.githubusercontent.com/ivolo/disposable-email-domains/master/index.json", &disposable); err != nil {
		log.Print(err)
		return
	}

	var disposableWild []string
	if err := getURLJson("https://raw.githubusercontent.com/ivolo/disposable-email-domains/master/wildcard.json", &disposableWild); err != nil {
		log.Print(err)
		return
	}

	freeProvider, err := getURLPHP("https://raw.githubusercontent.com/daveearley/Email-Validation-Tool/master/src/data/email-providers.php")
	if err != nil {
		log.Print(err)
		return
	}

	tlds, err := getURLStrings("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	if err != nil {
		log.Print(err)
		return
	}

	buf := &bytes.Buffer{}
	_, _ = fmt.Fprintln(buf, "package emailvalidator\n\n")
	printMap(buf, "disposableDomain", disposable)
	printMap(buf, "wildDisposableDomain", disposableWild)
	printMap(buf, "freeProvider", freeProvider)
	printMap(buf, "tlds", tlds)

	if err := formatCode(buf.Bytes(), fl); err != nil {
		log.Print(err)
		return
	}
}
