---
title: Coinbase
slug: /coinbase
---

## How do I import Coinbase data to Google Sheets?

Coinbase is an American company that operates a cryptocurrency exchange platform that operates remote-first without an official physical headquarters. The company was founded in 2012 by Brian Armstrong and Fred Ehrsam, and as of March 2021, was the largest cryptocurrency exchange in the United States by trading volume. In this guide, I'll walk you through how to import data from Coinbase into Google Sheets.

** BEFORE YOU BEGIN **

If you haven't already, install the Data Connector Add-on for Google Sheets from the [Google Workspace Marketplace](https://workspace.google.com/marketplace/app/appname/529655450076)

**Step 1: Create your API Request**

1. Open Google Sheets and click `Add-ons -> Data Connector -> Manage Connections`
2. Click `NEW COMMAND`
3. Name your command. In this case, we will name it `coinbase`
4. Select `API` for the Type
5. Select `GET` for the Method
6. Enter `https://api.coinbase.com/v2/exchange-rates?currency=BTC` in the URL field.
7. Select `JMESPath` for the Filter type
8. Enter `data.rates.USD` in the expression box.
9. Click `SAVE`

**Step 2: Run the command**

1. Now, you can run your command: `=run("coinbase")`