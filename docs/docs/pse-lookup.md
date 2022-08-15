---
title: PSE Lookup API
slug: /pse-lookup
---

## How do I import Philippine Stock Exchange data to Google Sheets?

PSE Lookup API was created by Vrymel Omandam and describes itself as 

> PSE Lookup is a personal project and not affiliated with the Philippine Stock Exchange.
The goal of this tool is to make downloading of Philippine Stock Exchange historical data readily available to the public. A project created out of frustration when searching for historical data myself.
> Please counter check the data you download as I can't guarantee for its accurateness. If you ever find something erroneous, please don't hesitate to message me at vrymel@gmail.com so I can address it.
> Lastly, you agree that PSE Lookup is not liable for any losses that may result from any data you downloaded from the tool.

Their website can be found at https://pselookup.vrymel.com/

** BEFORE YOU BEGIN **

If you haven't already, install the Data Connector Add-on for Google Sheets from the [Google Workspace Marketplace](https://workspace.google.com/marketplace/app/appname/529655450076)

**Step 1: Create your API Request**

1. Open Google Sheets and click `Add-ons -> Data Connector -> Manage Connections`
2. Click `NEW COMMAND`
3. Name your command. In this case, we will name it `pse`
4. Select `API` for the Type
5. Select `GET` for the Method
6. Enter `https://pselookup.vrymel.com/api/stocks/+++1+++` in the URL field.
7. Select `JMESPath` for the Filter type
8. Enter `price.close` in the expression box.
9. Click `SAVE`

**Step 2: Run the command**

1. Now, you can run your command: `=run("pse", "JFC")`