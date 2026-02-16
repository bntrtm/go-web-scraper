# go-web-scraper
go-based SEO analytics tool that reports on the internal linking profile of any website

## Installation

With go already installed, you can run:
`go install https://github.com/bntrtm/go-web-scraper@latest`

## Usage

```
go-web-scraper [options] <link>
```

A `link` *must* be provided to the program to start crawling from.

### Options
`Concurrency control: -c <int>`

An integer specifying the number of goroutines to use in parallel (default: 1).

`page limit: -p <int>`

An integer specifying the maximum number of pages to crawl before exiting the program (default: 3). When set to <= 0, there will be no limit, causing the crawler to crawl for every link until all are found.

