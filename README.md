# Exploring Go

This repository is a collection of small Go programs I’ve built while learning the language more deeply through hands-on problem solving.

Rather than treating Go as a purely academic exercise, I’m using this space to explore how it can be applied to real engineering tasks such as networking, concurrency, HTTP communication, and command-line tooling.

A lot of these projects are intentionally small, but each one helps me understand an important Go concept by building something practical.

## Why this repository exists

I created this repository to:
- strengthen my understanding of core Go concepts
- build small tools instead of only reading syntax examples
- document my progress in a way that is practical and portfolio-friendly
- explore Go from a systems and infrastructure perspective

## Projects

### Uptime Checker
A concurrent CLI tool that checks whether a list of URLs is reachable, resolves their IP addresses, and reports HTTP status codes.

**Concepts explored:**
- goroutines
- sync.WaitGroup
- HTTP clients
- timeouts
- DNS lookups
- URL parsing

## What I’m focusing on
Most of the projects in this repository are centered around:
- concurrency
- networking
- CLI applications
- backend/system-style problem solving

## Notes
This is an evolving repository. Some projects may remain intentionally simple, while others may be revisited and improved as my understanding grows.

## Inspiration

Project ideas in this repository were inspired by the excellent 
[Golang for DevOps](https://github.com/techiescamp/golang-for-devops) repository.

I used it as a guide for identifying Go exercises while 
implementing the solutions myself to better understand the language and its 
standard library.