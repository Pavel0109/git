package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

//type Molecul struct {
//	params map[string]string
//}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func main() {

	configurets := []string{"iupac_name", "average_molecular_weight", "status"}

	xmlFile, err := os.Open("serum_metabolites.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)

	//result:= new ([]map[string]string)
	//current := make(map[string]string)
	//var current  []string
	//var result [][]string
	var result_map []map[string]string
	result_map_length := 0
	for {
		//len := 0
		tok, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happend", tokenErr)
			break
		} else if tokenErr == io.EOF {
			break
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == "metabolite" {
				cs := make(map[string]string)
				result_map = append(result_map, cs)
				result_map_length++
			}
			if stringInSlice(tok.Name.Local, configurets) {
				var buf string
				if err := decoder.DecodeElement(&buf, &tok); err != nil {
					fmt.Println("error", err)
				}
				result_map[result_map_length-1][tok.Name.Local] = buf
				//fmt.Println(buf)
				//current = append(current,buf)

				//fmt.Println(len, result[len])

			}
			//case xml.EndElement:
			//	if tok.Name.Local == "metabolite" {
			//		result=append(result,current)
			//	}
			//}
		}
	}
	//fmt.Println(result_map)
	//for _, i := range result_map{
	//	fmt.Println(i)
	//}

	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(configurets)

	checkError("Cannot write to file", err)
	for _, current := range result_map {
		var text []string
		for _, value := range configurets {
			text = append(text, current[value])
		}
		fmt.Println(current)
		err := writer.Write(text)
		checkError("Cannot write to file", err)
	}
}
