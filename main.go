package main

import (
	"context"
	"samwang0723/jarvis/server/handlers"
)

func main() {
	handlers.BatchingDownload(context.Background(), -2, 5000)
}
