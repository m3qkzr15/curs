package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
)

type Rate struct {
	Currency string `xml:"currency,attr"`
	Value    string `xml:",chardata"`
}

type Cube struct {
	Rates []Rate `xml:"Rate"`
}

type Body struct {
	Cube Cube `xml:"Cube"`
}

type DataSet struct {
	XMLName xml.Name `xml:"DataSet"`
	Body    Body     `xml:"Body"`
}

func GetExchangeRate() float64 {

	url := "https://www.bnr.ro/nbrfxrates.xml"
	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var data DataSet
	err = xml.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	cursEur, err := strconv.ParseFloat(data.Body.Cube.Rates[10].Value, 64)
	if err != nil {
		panic(err)
	}

	return cursEur

}
