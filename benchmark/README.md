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

- Updated (UTC): 2026-03-05T07:03:38Z
- Go: `go1.26.0`
- OS/Arch: `linux/amd64`
- CPU: `AMD Ryzen 7 5800X 8-Core Processor`

### Registration

| Library |     ops (N) |     ns/op |      B/op | allocs/op |
|---------|------------:|----------:|----------:|----------:|
| Dig     |      48,283 |    24,177 |    24,778 |       349 |
| Do      |     179,978 |     6,612 |     4,539 |        74 |
| Remy    | **839,143** | **1,407** | **1,104** |    **33** |

### Singleton Retrieval

| Library |        ops (N) |     ns/op |  B/op | allocs/op |
|---------|---------------:|----------:|------:|----------:|
| Dig     |      1,374,896 |       876 |   576 |        20 |
| Do      |      4,580,288 |     263.4 |   176 |         5 |
| Remy    | **41,592,746** | **27.99** | **0** |     **0** |

### Factory Retrieval

| Library |       ops (N) |     ns/op |    B/op | allocs/op |
|---------|--------------:|----------:|--------:|----------:|
| Dig     |     1,510,795 |     806.2 |     576 |        20 |
| Do      |       191,050 |     6,014 |   3,016 |        67 |
| Remy    | **2,585,070** | **470.4** | **120** |     **4** |

### Instance Retrieval

| Library |        ops (N) |     ns/op |  B/op | allocs/op |
|---------|---------------:|----------:|------:|----------:|
| Dig     |      1,521,675 |       793 |   576 |        20 |
| Do      |      4,793,290 |       250 |   176 |         5 |
| Remy    | **39,992,311** | **28.48** | **0** |     **0** |

### Nested Dependency Resolution

| Library |       ops (N) |     ns/op |  B/op | allocs/op |
|---------|--------------:|----------:|------:|----------:|
| Dig     |       269,406 |     4,255 | 3,072 |       104 |
| Do      |       722,618 |     1,443 |   960 |        26 |
| Remy    | **8,763,864** | **134.2** | **0** |     **0** |

### Multiple Retrievals

| Library |       ops (N) |     ns/op |  B/op | allocs/op |
|---------|--------------:|----------:|------:|----------:|
| Dig     |       297,512 |     4,019 | 2,880 |       100 |
| Do      |       826,291 |     1,339 |   880 |        25 |
| Remy    | **7,974,126** | **143.5** | **0** |     **0** |

### Unregistered Type Retrieval

| Library |       ops (N) |     ns/op |   B/op | allocs/op |
|---------|--------------:|----------:|-------:|----------:|
| Dig     |       118,795 |    10,311 |  5,529 |       151 |
| Do      |       138,446 |     8,447 | 10,210 |       134 |
| Remy    | **5,189,788** | **232.4** | **80** |     **5** |

<!-- BENCHSTAT:END -->

## Benchstat Output

- Raw file: [`.benchstat/benchstat.txt`](./.benchstat/benchstat.txt)
