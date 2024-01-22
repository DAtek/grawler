[![codecov](https://codecov.io/gh/DAtek/grawler/graph/badge.svg?token=3WMKQDRJ95)](https://codecov.io/gh/DAtek/grawler) [![Go Report Card](https://goreportcard.com/badge/github.com/DAtek/grawler)](https://goreportcard.com/report/github.com/DAtek/grawler)

# Grawler

## Simple and performant web crawler in Go

<img src="./gopher.png" width="200" />

### How it works
The crawler is using 2 types of workers:
- **Page loaders**
- **Page analyzers**

**Page loaders** are consuming the **remaining URL channel** and are downloading pages from the internet and putting them into a **cache**, also putting the downloaded page's URL into the **downloaded URL channel**.

**Page analyzers** are consuming the **downloaded URL channel** and reading the page's content from the **cache**, then analyzing the content, extracting additional URLs and the wanted model (if possible). The extracted new URLs are being put into the **remaining URL channel**, the found model in the **result channel**.

The whole process is being started with putting the starting URL into the **remaining URL channel**.

The number of **Page loaders** and **Page analyzers** are configurable.

Your possibilities are endless: you can implement your own **cache**, **page loader** and **analyzer**, the mocks and interfaces in the source will help you.

For guidance, please have a look at `crawler_test.go`.

The gopher was made with the [Gopher Konstructor](https://quasilyte.dev/gopherkon)
