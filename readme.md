# k6 Reporter Extension

- simple extension to output test results as html, to be used with handleSummary, based on https://github.com/benc-uk/k6-reporter
- on your test script use:

`import reporter from 'k6/x/reporter';`

and invoke `generateReport` as the output method for handleSummary

`reporter.generateReport(JSON.stringify(data), "Report Title")`