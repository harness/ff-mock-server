package config

// Options holds cli flags
var Options struct {
	Timeout        *int     `short:"t" long:"timeout" description:"Request timeout"`
	StatusCode     *int     `short:"s" long:"status-code" description:"returns HTTP status code"`
	Message        string   `short:"m" long:"message" description:"Message to display in response"`
	SSEOffSequence []int    `short:"e" long:"sse" description:"SSEOffSequence off sequence in sec"`
	SSEOffDuration *int     `long:"sse-out" description:"SSEOffSequence off time in sec"`
	Handlers       []string `short:"o" long:"operation" description:"operation"`
}
