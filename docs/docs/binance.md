---
title: Binance
slug: /binance
---

## How do I import Binance data to Google Sheets?

Binance is a cryptocurrency exchange that provides a platform for trading various cryptocurrencies. As of April 2021, Binance was the largest cryptocurrency exchange in the world in terms of trading volume. Binance was founded by Changpeng Zhao, a developer who had previously created high frequency trading software. Binance was initially based in China, but later moved out of China due to China's increasing regulation of cryptocurrency. Binance describes themselves as "Beyond operating the world's leading cryptocurrency exchange, Binance spans an entire ecosystem". In this guide, I'll walk you through how to import data from Binance into Google Sheets.

** BEFORE YOU BEGIN **

If you haven't already, install the Data Connector Add-on for Google Sheets from the [Google Workspace Marketplace](https://workspace.google.com/marketplace/app/appname/529655450076)

**Step 1: Create your API Request**

1. Open Google Sheets and click `Add-ons -> Data Connector -> Manage Connections`
2. Click `NEW COMMAND`
3. Name your command. In this case, we will name it `binance`
4. Select `API` for the Type
5. Select `GET` for the Method
6. Enter `https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT` in the URL field.
7. Select `JMESPath` for the Filter type
8. Enter `price` in the expression box.
9. Click `SAVE`

**Step 2: Run the command**

1. Now, you can run your command: `=run("binance")`