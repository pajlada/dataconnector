---
title: Stripe
slug: /stripe
---

## How do I import Stripe data to Google Sheets?

Stripe is one of the largest payment processors for the web. Stripe’s products power payments for online and in-person retailers, subscriptions businesses, software platforms and marketplaces, and everything in between. In this guide, I'll walk you through how to import data from Stripe into Google Sheets.

** BEFORE YOU BEGIN **

If you haven't already, install the Data Connector Add-on for Google Sheets from the [Google Workspace Marketplace](https://workspace.google.com/marketplace/app/appname/529655450076)

**Step 1: Get your Stripe Secret key**

Login to the [Stripe Dashboard](https://dashboard.stripe.com/dashboard) and click on the `Get your API keys` link. 

![Stripe API Key](/img/stripe/stripe1.png)

Click on the button to reveal your secret key. You will need this in Step 3 below.

**Step 2: Create your API Request**

1. Open Google Sheets and click `Add-ons -> Data Connector -> Manage Connections`
2. Click `NEW COMMAND`
3. Name your command. In this case, we will name it `stripe`
4. Select `API` for the Type
5. Select `GET` for the Method
6. Enter `https://api.stripe.com/v1/charges?limit=10` in the URL field.
7. In the Headers section, enter `Authorization` for the Key and `Bearer +++1+++` for the Value.
8. Select `JMESPath` for the Filter type
9. Enter `data[*].[amount]` in the expression box.
10. Click `SAVE`

**Step 3: Run the command**

1. Enter your Stripe Secret Key in any cell. In our case, we chose cell `B1`.
2. Now, you can run your command: `=run("stripe", B1)`