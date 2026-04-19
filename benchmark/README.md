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

- Updated (UTC): 2026-04-19T06:05:42Z
- Go: `go1.26.2`
- OS/Arch: `linux/amd64`
- CPU: `AMD EPYC 7763 64-Core Processor`

### Registration

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 37,849 | 31,733 | 24,764 | 349 |
| Do | 239,636 | 4,803 | 3,376 | 54 |
| Gontainer | 100,720 | 11,724 | 5,777 | 134 |
| Remy | **659,592** | **1,817** | **1,104** | **33** |

### Singleton Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 1,000,000 | 1,027 | 576 | 20 |
| Do | 3,675,618 | 321.7 | 176 | 5 |
| Gontainer | 3,739,525 | 321.8 | 48 | 3 |
| Remy | **28,633,857** | **41.7** | **0** | **0** |

### Factory Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 1,000,000 | 1,027 | 576 | 20 |
| Do | 149,560 | 7,825 | 3,016 | 67 |
| Gontainer | 946,182 | 1,218 | 280 | 11 |
| Remy | **1,838,446** | **640.9** | **120** | **4** |

### Instance Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 1,000,000 | 1,025 | 576 | 20 |
| Do | 3,790,702 | 314.9 | 176 | 5 |
| Gontainer | 3,527,917 | 340.3 | 48 | 3 |
| Remy | **28,641,868** | **41.57** | **0** | **0** |

### Nested Dependency Resolution

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 213,824 | 5,531 | 3,072 | 104 |
| Do | 606,861 | 1,791 | 960 | 26 |
| Gontainer | 325,532 | 3,581 | 792 | 36 |
| Remy | **6,077,904** | **198.2** | **0** | **0** |

### Multiple Retrievals

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 227,192 | 5,195 | 2,880 | 100 |
| Do | 674,749 | 1,622 | 880 | 25 |
| Gontainer | 201,486 | 5,922 | 1,400 | 55 |
| Remy | **5,698,176** | **210.3** | **0** | **0** |

### Unregistered Type Retrieval

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| Dig | 86,864 | 13,676 | 5,526 | 151 |
| Do | 108,690 | 11,125 | 10,203 | 134 |
| Gontainer | 642,919 | 1,770 | 384 | 10 |
| Remy | **3,978,879** | **293.2** | **80** | **5** |

<!-- BENCHSTAT:END -->

## Benchstat Output

- Raw file: [`.benchstat/benchstat.txt`](./.benchstat/benchstat.txt)
