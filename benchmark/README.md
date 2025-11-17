# Benchmark Results

## How to Run

```sh
# Generate benchmark setup
go generate ./...

# Run benchmarks with memory allocation info
go test -bench=. -benchmem ./benchmark/tests
```

---

## Results

All benchmark results below were obtained on the following test environment:

- **OS:** Linux
- **Architecture:** amd64
- **CPU:** AMD Ryzen 7 5800X (8-Core Processor)

### Registration

| Benchmark | ops (N) | ns/op  | B/op   | allocs/op |
|-----------|---------|--------|--------|-----------|
| **Dig**   | 49,867  | 24,437 | 24,641 | 332       |
| **Do**    | 303,548 | 3,735  | 3,378  | 54        |
| **Remy**  | 683,485 | 1,742  | 1,208  | 35        |

### Singleton Retrieval

| Benchmark | ops (N)        | ns/op     | B/op  | allocs/op |
|-----------|----------------|-----------|-------|-----------|
| **Dig**   | 1,628,277      | 724.9     | 544   | 16        |
| **Do**    | 4,998,186      | 237.4     | 176   | 5         |
| **Remy**  | **24,148,269** | **47.40** | **0** | **0**     |

### Factory Retrieval

| Benchmark | ops (N)   | ns/op | B/op  | allocs/op |
|-----------|-----------|-------|-------|-----------|
| **Dig**   | 1,638,084 | 730.1 | 544   | 16        |
| **Do**    | 204,300   | 5,707 | 3,032 | 69        |
| **Remy**  | 1,786,933 | 669.1 | 136   | 6         |

### Instance Retrieval

| Benchmark | ops (N)        | ns/op     | B/op  | allocs/op |
|-----------|----------------|-----------|-------|-----------|
| **Dig**   | 1,656,339      | 728.2     | 544   | 16        |
| **Do**    | 5,117,502      | 234.4     | 176   | 5         |
| **Remy**  | **25,066,693** | **47.48** | **0** | **0**     |

### Nested Dependency Resolution

| Benchmark | ops (N)       | ns/op     | B/op  | allocs/op |
|-----------|---------------|-----------|-------|-----------|
| **Dig**   | 310,790       | 3,919     | 2,912 | 84        |
| **Do**    | 847,618       | 1,300     | 960   | 26        |
| **Remy**  | **5,563,692** | **211.3** | **0** | **0**     |

### Multiple Retrievals

| Benchmark | ops (N)       | ns/op     | B/op  | allocs/op |
|-----------|---------------|-----------|-------|-----------|
| **Dig**   | 301,608       | 3,682     | 2,720 | 80        |
| **Do**    | 1,013,702     | 1,195     | 880   | 25        |
| **Remy**  | **5,083,998** | **237.0** | **0** | **0**     |
