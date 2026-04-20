# Benchmark Results

## How to Run

```sh
# Install benchmark tool
make -C benchmark install-tools

# Run all three benchmark suites, run benchstat, and auto-update this README
make -C benchmark bench

# Exact benchstat command used for the final comparison
cd benchmark/.benchstat && benchstat dig.txt do.txt gontainer.txt remy.txt
```

## Notes

- Raw benchmark outputs are generated in `benchmark/.benchstat/`.
- The section below is automatically updated by `make -C benchmark bench` and CI.

<!-- BENCHSTAT:START -->

## Latest Benchmark Comparison

- Updated (UTC): 2026-04-20T06:43:45Z
- Go: `go1.26.2`
- OS/Arch: `linux/amd64`
- CPU: `AMD EPYC 9V74 80-Core Processor`

### Registration

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 37,299 | 32,401 | 24,764 | 349 |
| Do | 248,014 | 4,827 | 3,376 | 54 |
| Gontainer | 108,736 | 10,924 | 5,777 | 134 |
| Remy | **667,326** | **1,663** | **1,104** | **33** |

### Singleton Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 1,260,001 | 953.9 | 576 | 20 |
| Do | 4,175,308 | 287.6 | 176 | 5 |
| Gontainer | 3,931,974 | 306.5 | 48 | 3 |
| Remy | **30,125,343** | **39.91** | **0** | **0** |

### Factory Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 1,260,044 | 945.7 | 576 | 20 |
| Do | 155,140 | 7,526 | 3,016 | 67 |
| Gontainer | 964,610 | 1,111 | 280 | 11 |
| Remy | **1,972,420** | **611.7** | **120** | **4** |

### Instance Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 1,268,157 | 955.6 | 576 | 20 |
| Do | 4,267,882 | 277.6 | 176 | 5 |
| Gontainer | 3,661,696 | 324.6 | 48 | 3 |
| Remy | **30,568,904** | **39.26** | **0** | **0** |

### Nested Dependency Resolution

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 221,380 | 5,096 | 3,072 | 104 |
| Do | 660,878 | 1,586 | 960 | 26 |
| Gontainer | 339,559 | 3,349 | 792 | 36 |
| Remy | **6,536,043** | **184** | **0** | **0** |

### Multiple Retrievals

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 237,930 | 4,851 | 2,880 | 100 |
| Do | 731,650 | 1,433 | 880 | 25 |
| Gontainer | 209,180 | 5,573 | 1,400 | 55 |
| Remy | **6,034,513** | **198.8** | **0** | **0** |

### Unregistered Type Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 90,512 | 12,969 | 5,527 | 151 |
| Do | 117,328 | 10,208 | 10,203 | 134 |
| Gontainer | 704,574 | 1,501 | 384 | 10 |
| Remy | **4,333,286** | **274.3** | **80** | **5** |

<!-- BENCHSTAT:END -->

## Benchstat Output

- Raw file: [`.benchstat/benchstat.txt`](./.benchstat/benchstat.txt)
