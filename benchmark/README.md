# Benchmark Results

## How to Run

```sh
# Install benchmark tool
make -C benchmark install-tools

# Run all three benchmark suites, run benchstat, and auto-update this README
make -C benchmark bench

# Exact benchstat command used for the final comparison
cd benchmark/.benchstat && benchstat dig.txt do.txt remy.txt
```

## Notes

- Raw benchmark outputs are generated in `benchmark/.benchstat/`.
- The section below is automatically updated by `make -C benchmark bench` and CI.

<!-- BENCHSTAT:START -->

## Latest Benchmark Comparison

- Updated (UTC): 2026-03-05T08:07:48Z
- Go: `go1.26.0`
- OS/Arch: `linux/amd64`
- CPU: `AMD EPYC 9V74 80-Core Processor`

### Registration

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 37,621 | 32,211 | 24,764 | 349 |
| Do | 250,491 | 4,684 | 3,376 | 54 |
| Remy | **697,466** | **1,596** | **1,104** | **33** |

### Singleton Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 1,259,188 | 951.4 | 576 | 20 |
| Do | 4,210,194 | 286.7 | 176 | 5 |
| Remy | **32,114,666** | **37.4** | **0** | **0** |

### Factory Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 1,261,896 | 961.1 | 576 | 20 |
| Do | 148,128 | 7,892 | 3,016 | 67 |
| Remy | **1,990,269** | **604.4** | **120** | **4** |

### Instance Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 1,247,976 | 959.2 | 576 | 20 |
| Do | 4,350,697 | 273.2 | 176 | 5 |
| Remy | **32,007,944** | **35.99** | **0** | **0** |

### Nested Dependency Resolution

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 228,612 | 5,115 | 3,072 | 104 |
| Do | 662,799 | 1,586 | 960 | 26 |
| Remy | **6,945,984** | **174.3** | **0** | **0** |

### Multiple Retrievals

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 241,912 | 4,826 | 2,880 | 100 |
| Do | 696,375 | 1,438 | 880 | 25 |
| Remy | **6,789,886** | **188.7** | **0** | **0** |

### Unregistered Type Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 91,122 | 12,954 | 5,526 | 151 |
| Do | 115,471 | 10,068 | 10,203 | 134 |
| Remy | **4,332,814** | **274.5** | **80** | **5** |

<!-- BENCHSTAT:END -->

## Benchstat Output

- Raw file: [`.benchstat/benchstat.txt`](./.benchstat/benchstat.txt)
