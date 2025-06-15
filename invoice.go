package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createInvoiceList(invoiceList *fyne.Container) error {

	storedInvoices, err := getAllRedisKeys(redisDB)
	if err != nil {
		panic(err)
	}

	for DBkey, value := range storedInvoices {

		var inv Invoice

		err := json.Unmarshal([]byte(value), &inv)
		if err != nil {
			panic(err)
		}

		rowId := widget.NewLabel(fmt.Sprintf("Id %v", DBkey))
		rowLabel := widget.NewLabel(fmt.Sprintf("Factura %v", inv.CompanyName))
		rowButtonDelete := widget.NewButton("Sterge", func() {
			deleteInvoice(invoiceList, DBkey)
		})

		rowButtonDetails := widget.NewButton("Detalii", func() {
			detailsInvoice(inv)
		})

		row := container.NewHBox(rowId, rowLabel, layout.NewSpacer(), rowButtonDetails, rowButtonDelete)
		invoiceList.Add(row)
	}

	return nil
}

func deleteInvoice(invoiceList *fyne.Container, dbkey string) {

	err := deleteRedisKey(redisDB, dbkey)
	if err != nil {
		panic(err)
	}

	invoiceList.Objects = nil

	err = createInvoiceList(invoiceList)
	if err != nil {
		panic(err)
	}

}

func addInvoice(invoiceList *fyne.Container, companyName, companyCUI, product, price, quantity string) error {

	if !isFloat(price) {
		err := errors.New("pretul nu este corect")
		return err
	}

	if !isPositiveInt(quantity) {
		err := errors.New("cantitatea nu este corecta")
		return err
	}

	if !lenBetween(companyName, 3, 50) {
		err := errors.New("numele companiei nu este corect")
		return err
	}

	if !lenBetween(companyCUI, 1, 20) {
		err := errors.New("cuiul companiei nu este corect")
		return err
	}

	if product == "" {
		err := errors.New("trebuie selectat un produs")
		return err
	}

	invoiceCount += 1

	convertedPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return err
	}
	totalPrice := convertedPrice * cursBNR

	newInvoice := Invoice{
		ID:          invoiceCount,
		CompanyName: companyName,
		CompanyCUI:  companyCUI,
		Product:     product,
		PriceEur:    convertedPrice,
		PriceRon:    totalPrice,
		Quantity:    quantity,
	}

	json, err := json.MarshalIndent(newInvoice, "", " ")
	if err != nil {
		return err
	}

	err = setRedisKey(redisDB, strconv.Itoa(invoiceCount), string(json))
	if err != nil {
		return err
	}

	invoiceList.Objects = nil

	err = createInvoiceList(invoiceList)
	if err != nil {
		return err
	}

	invoiceList.Refresh()

	return nil
}

func isFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isPositiveInt(s string) bool {
	v, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return v > 0
}

func lenBetween(s string, min, max int) bool {
	if min > max || min < 0 {
		return false
	}
	n := utf8.RuneCountInString(s)
	return n >= min && n <= max
}

type Invoice struct {
	ID          int     `json:"id"`
	CompanyName string  `json:"companyName"`
	CompanyCUI  string  `json:"companyCUI"`
	Product     string  `json:"product"`
	PriceEur    float64 `json:"priceEur"`
	PriceRon    float64 `json:"priceRon"`
	Quantity    string  `json:"quantity"`
}
