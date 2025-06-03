package main

import (
    _ "./heuristics"
)

func main() {
    StartProxy(":8080", "http://localhost:8081")
}