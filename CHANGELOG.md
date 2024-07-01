# Changelog

All notable changes to this project will be documented in this file.

## [unreleased]

### Bug Fixes

- Improve performance while calculating average
- Adjust pulling count and limit
- Add more information and check more conditions
- Add linter and adjust makefile
- Adjust lint errors
- Adjust lint errors
- Update makefile to solve gofumpt skip problems
- Update makefile and docker build version
- Remove analysis into database, don't prefer to put calculation result in
- Finetune performance of doing selection
- Add stock market indicator
- Update html layout
- Adjust after close quote changes
- Adjust html format
- Adjust html format
- Adjust html format
- Adjust html format
- Testing code review ([#26](https://github.com/samwang0723/stock-crawler/issues/26))
- Javascript math.Max Infinity issue fix
- Add comments and update halfyear threshold
- Make link to create new tab and adjust daily high selection
- Add missing feature ([#30](https://github.com/samwang0723/stock-crawler/issues/30))
- Add viewport to have mobile responsive size
- Add web icons
- Update javascript domain connection
- Adjust web layout
- Make the candlestick chart to be pixel perfect
- Html errors
- Solve the historical analysis using wrong date price issue
- Update the strategy of stock selection and correct some issues
- Make sure kafka always connected from consumer
- Remove sensitive info
- Use sentinel mode and remove the istio gateway for 6379 port
- Using correct password key for redis
- Support CORS
- Use password on sentinel
- Always re-pulling image
- Correct quoteChange %
- Support redis distributed lock for multiple pod insance
- Debugging redis sentinel
- Update redis distributed lock
- Update redis sentinel master issue
- Using https and refine redis query
- Support live/ready healthcheck and swagger
- Remove swagger-ui, only keep json api document
- Add more details in events
- Protect redis unmarshal error and correct closed order profit % calculation
- Solve profit calculation errors
- Solve add self-picked without authentication
- Clear cookie if login failed
- Deprecate webscraping and have weekly close table
- Improve performance by removing log and unused feature ([#41](https://github.com/samwang0723/stock-crawler/issues/41))
- Use rs/cors package to handle cors
- Use rs/cors package to handle cors
- Update swagger to have security bearer param
- Respond with 401 when login failed
- Solve add self-picked duplication problem
- Solve NULL duplication problem
- Add ingress to handle frontend server
- Support volume check on selection
- Adjust average volume calc logic
- Solve unrealized profit issue and update selection criteria
- Add datadog redis metrics
- Update prometheus monitoring on MySQL
- Update MySQL auth method and reduce warning logs
- Remove small memory limit
- Upgrade go version to 1.22 and remove swagger ([#43](https://github.com/samwang0723/stock-crawler/issues/43))

### Features

- Add average data into daily closes ([#24](https://github.com/samwang0723/stock-crawler/issues/24))
- Insert average into every daily_closes records
- Add selection for stock picking
- Add analysis
- Add analysis into selection api
- Add github workflow support
- Update the algorithm of selection
- Support realtime price monitoring based on yesterday slections ([#25](https://github.com/samwang0723/stock-crawler/issues/25))
- Prevent pulling from public holiday
- Update more detailed information
- Support picked stock selection ([#28](https://github.com/samwang0723/stock-crawler/issues/28))
- Expose realtime data in self picked stock list
- Support mysql helm chart with istio-gateway
- Support helm chart deploy redis sentinel
- Support kafka and deploy app
- Support user login and with transaction / balance ([#31](https://github.com/samwang0723/stock-crawler/issues/31))
- Support local environment without crawling real web
- Support crawling realtime price for open orders ([#38](https://github.com/samwang0723/stock-crawler/issues/38))
- Support order viewing page
- Support login method ([#39](https://github.com/samwang0723/stock-crawler/issues/39))
- Update for operation
- Change to smartproxy
- Support realtime date
- Update indexes for performance
- Support grafana monitoring
- Update prometheus redis monitoring
- Changed to Kong gateway

### Miscellaneous Tasks

- Adding cloudflare worker script to cache API result

### Refactor

- Fix domain:port and change the chart
- Using the official eventsourcing way ([#32](https://github.com/samwang0723/stock-crawler/issues/32))
- Refined to all using zerolog ([#35](https://github.com/samwang0723/stock-crawler/issues/35))
- Refine config usage and introduce sentry.io ([#36](https://github.com/samwang0723/stock-crawler/issues/36))
- Add order ([#37](https://github.com/samwang0723/stock-crawler/issues/37))
- Refine order page style
- Refine UX style

### Build

- Bump golang.org/x/net from 0.5.0 to 0.7.0 ([#29](https://github.com/samwang0723/stock-crawler/issues/29))
- Bump golang.org/x/net from 0.7.0 to 0.17.0 ([#33](https://github.com/samwang0723/stock-crawler/issues/33))

<!-- generated by git-cliff -->
