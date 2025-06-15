package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var cursApp fyne.App
var cursBNR float64
var redisDB = createRedisClient()
var invoiceCount = 0

func main() {

	cursApp = app.New()
	appWindow := cursApp.NewWindow("Curs BNR App")
	appWindow.Resize(fyne.NewSize(500, 350))
	appWindow.SetContent(container.NewVBox(makeUI()))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nApp is shutting down, flushing Redis DB...")

		if err := redisDB.FlushDB(ctx).Err(); err != nil {
			fmt.Println("Error flushing Redis DB:", err)
		} else {
			fmt.Println("Redis DB flushed successfully.")
		}

		os.Exit(0)
	}()

	appWindow.SetOnClosed(func() {
		if err := redisDB.FlushDB(ctx).Err(); err != nil {
			fmt.Println("Error flushing Redis DB:", err)
		} else {
			fmt.Println("Redis DB flushed successfully.")
		}
	})

	appWindow.Show()
	cursApp.Run()
}

func makeUI() fyne.CanvasObject {

	cursBNR = GetExchangeRate()
	labelText := fmt.Sprintf("Cursul BNR curent: %.4f", cursBNR)
	cursBNROutput := widget.NewLabel(labelText)

	invoiceListLabel := widget.NewLabel("Lista Facturilor emise")
	invoiceList := container.NewVBox()

	err := createInvoiceList(invoiceList)
	if err != nil {
		panic(err)
	}

	invoiceListScroll := container.NewScroll(invoiceList)
	invoiceListScroll.SetMinSize(fyne.NewSize(350, 200))

	addInvoiceButton := widget.NewButton("Creeaza Factura", func() {
		newInvoice(invoiceList)
	})

	content := container.NewVBox(cursBNROutput, invoiceListLabel, invoiceListScroll, addInvoiceButton)

	return content
}

func newInvoice(invoiceList *fyne.Container) {

	addInvoiceWindow := cursApp.NewWindow("Creeaza Factura")
	addInvoiceWindow.Resize(fyne.NewSize(300, 300))

	label := widget.NewLabel("Completeaza Detalii")

	companyName := widget.NewEntry()
	companyName.SetPlaceHolder("Nume Companie...")

	companyCUI := widget.NewEntry()
	companyCUI.SetPlaceHolder("CUI Companie...")

	product := widget.NewSelect([]string{"Produs I", "Produs II", "Produs III"}, func(s string) {})
	product.PlaceHolder = "Alege produs..."

	price := widget.NewEntry()
	price.SetPlaceHolder("Price EUR.")

	quantity := widget.NewEntry()
	quantity.SetPlaceHolder("Cant.")

	submitButton := widget.NewButton("Creeaza Factura", func() {
		err := addInvoice(invoiceList, companyName.Text, companyCUI.Text, product.Selected, price.Text, quantity.Text)
		if err != nil {
			dialog.ShowError(err, addInvoiceWindow)
		} else {
			addInvoiceWindow.Close()
		}

	})

	newInvoiceForm := container.NewVBox(label, companyName, companyCUI, product, price, quantity, submitButton)

	addInvoiceWindow.SetContent(newInvoiceForm)

	addInvoiceWindow.Show()

}

func detailsInvoice(inv Invoice) {

	title := fmt.Sprintf("Detalii Factura %v", inv.ID)
	detailsWindow := cursApp.NewWindow(title)

	detailsWindow.Resize(fyne.NewSize(300, 500))

	companyName := widget.NewLabel("Nume companie")
	companyName.TextStyle = fyne.TextStyle{Bold: true}
	companyNameValue := widget.NewLabel(inv.CompanyName)

	companyCUI := widget.NewLabel("CUI companie")
	companyCUI.TextStyle = fyne.TextStyle{Bold: true}
	companyCUIValue := widget.NewLabel(inv.CompanyCUI)

	product := widget.NewLabel("Produs")
	product.TextStyle = fyne.TextStyle{Bold: true}
	productValue := widget.NewLabel(inv.Product)

	priceRon := widget.NewLabel("Pret RON")
	priceRon.TextStyle = fyne.TextStyle{Bold: true}
	priceRonValue := widget.NewLabel(strconv.FormatFloat(inv.PriceRon, 'f', 2, 64))

	priceEur := widget.NewLabel("Pret EUR")
	priceEur.TextStyle = fyne.TextStyle{Bold: true}
	priceEurValue := widget.NewLabel(strconv.FormatFloat(inv.PriceEur, 'f', 2, 64))

	quantity := widget.NewLabel("Cantitate")
	quantity.TextStyle = fyne.TextStyle{Bold: true}
	quantityValue := widget.NewLabel(inv.Quantity)

	detailsInvoiceBox := container.NewVBox(companyName, companyNameValue, companyCUI, companyCUIValue, product, productValue, priceRon, priceRonValue, priceEur, priceEurValue, quantity, quantityValue)

	detailsWindow.SetContent(detailsInvoiceBox)

	detailsWindow.Show()
}
