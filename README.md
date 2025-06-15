# Curs BNR Invoice App

This is a simple desktop application written in Go using the [Fyne GUI toolkit](https://fyne.io/). It interacts with Redis to temporarily store invoice data and shows the current BNR (Romanian National Bank) exchange rate.

## 🧾 Features

- Displays the current BNR exchange rate.
- Allows users to create invoices with:
  - Company name
  - Company CUI
  - Product selection
  - Price in EUR
  - Quantity
- Shows a list of all invoices created in the session.
- Opens detailed views for each invoice.
- Uses Redis to store invoice data during the app session.
- Flushes the Redis database automatically when the app closes.

## 🧰 Requirements

- Go 1.20+
- Redis running locally (default: `localhost:6379`)
- Fyne v2 GUI toolkit  
  ```bash
  go get fyne.io/fyne/v2