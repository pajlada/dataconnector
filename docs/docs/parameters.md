---
title: Parameters
slug: /parameters
---

## How do I parameterize commands?

In the `Getting Started` section, we created a command with a parameter: `+++1+++`. We could have easily hardcoded the parameter instead: `http://api.icndb.com/jokes/random/3`. Parameters, however, let you reuse your queries and are great if you are retrieving, say, stock data for multiple tickers. Without them, you'd have to create a new request for each ticker. Boo!

You can also add multiple parameters to a request, and they can be in the headers, body, and other parts of the request. To parameterize other parts of the request, just keep incrementing - e.g. the next parameter would be `+++2+++`, then `+++3+++` and so forth. In your formula, you can pass in parameters in one of 2 ways:

1. A comma-separated string like this: `=run("chuck","3,param2,param3",$A$1)`
2. A cell reference like this: `=run("chuck",C1:E1,$A$1)`, where cells `C1:E1` contain the values you want to pass in.
