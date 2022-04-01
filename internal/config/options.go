package config

var Options struct {
	Timeout    *int     `short:"t" long:"timeout" description:"Request timeout"`
	StatusCode *int     `short:"s" long:"status-code" description:"returns HTTP status code"`
	Message    string   `short:"m" long:"message" description:"Message to display in response"`
	SSE        []int    `short:"e" long:"sse" description:"SSE off sequence"`
	Handlers   []string `short:"o" long:"operation" description:"operation"`
}
