# redis-timeout

A simple Go program to demonstrate Redis context timeout behavior.

## Purpose

This program continuously reads a key from Redis in a loop, using a 5-second context timeout for each request. It helps illustrate the difference between:

- **Connection refused** — When Redis is stopped/killed (immediate error)
- **Context deadline exceeded** — When Redis is unresponsive but connection isn't refused (timeout after 5 seconds)

## Prerequisites

- Go 1.16+
- Docker

## Usage

### 1. Start Redis in Docker

```bash
docker run -d --name redis -p 6379:6379 redis
```

### 2. Run the program

```bash
go run main.go
```

You should see output like:

```
Connected to Redis: PONG
[14:30:01] Key 'mykey' does not exist
[14:30:02] Key 'mykey' does not exist
...
```

### 3. Pause Redis to trigger context deadline

In a separate terminal, pause the Redis container:

```bash
docker pause redis
```

### 4. Observe the context deadline error

After pausing Redis, wait ~5 seconds. You'll see the timeout error:

```
[14:30:15] Error reading key: context deadline exceeded
[14:30:21] Error reading key: context deadline exceeded
...
```

Each error takes ~5 seconds to appear because the program waits for the context timeout.

### 5. Unpause Redis to resume normal operation

```bash
docker unpause redis
```

The program will immediately start working again:

```
[14:30:25] Key 'mykey' does not exist
...
```

## Cleanup

```bash
docker stop redis
docker rm redis
```

## Note

If you **stop** Redis (`docker stop redis`) instead of pausing it, you'll get a `connection refused` error immediately — not a timeout. This is because stopping Redis closes the port, and the OS rejects the connection instantly before any timeout can occur.
