---
title: Getting Started
slug: /
---

## Step 1: Install the Add-on

If you haven't already, install the Data Connector Add-on for Google Sheets from the [Google Workspace Marketplace](https://workspace.google.com/marketplace/app/appname/529655450076)

## Step 2: Run your first command

Open the Data Connector side menu and click the "New Command" button

1. Name your command `chuck`
2. Select `API` for the `Type`
3. Select `GET` for the `Method`
4. Enter `http://api.icndb.com/jokes/random/+++1+++` in the `URL` field
4. Select `JMESPath` as the `Filter type`
5. Enter `value[*].[joke]` in the `Expression` field
6. Click `SAVE`

Now you're ready to run your command:

1. Enter `=run("chuck","3",$A$1)` in cell `B1`
2. You'll see a list of 3 jokes. Enter a different number to get more or less jokes
3. To refresh the formula, change the value in cell `A1`

** A quick note on parameters **

In the above command, we added a parameter: `+++1+++`. We could have easily hardcoded the parameter instead: `http://api.icndb.com/jokes/random/3`. Parameters, however, let you reuse your queries and are great if you are retrieving, say, stock data for multiple tickers. Without them, you'd have to create a new request for each ticker. Boo!

You can also add multiple parameters to a request, and they can be in the headers, body, and other parts of the request. To parameterize other parts of the request, just keep incrementing - e.g. the next parameter would be `+++2+++`, then `+++3+++` and so forth. In your formula, you can pass in parameters in one of 2 ways:

1. A comma-separated string like this: `=run("chuck","3,param2,param3",$A$1)`
2. A cell reference like this: `=run("chuck",C1:E1,$A$1)`, where cells `C1:E1` contain the values you want to pass in.

## That's it!

Congratulations, you just ran your first parameterized query! If you have any questions or need to connect to a specific data source, please [open an issue](https://github.com/dataconnector) and we'll add it here.
