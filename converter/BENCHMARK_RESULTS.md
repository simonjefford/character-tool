# Regex Optimization Benchmark Results

## Summary

Moving regular expression compilation from function-level to package-level variables resulted in significant performance improvements across all benchmarked functions.

## Performance Improvements

| Benchmark | Before (ns/op) | After (ns/op) | Speedup | Memory Before (B/op) | Memory After (B/op) | Allocs Before | Allocs After |
|-----------|----------------|---------------|---------|----------------------|---------------------|---------------|--------------|
| ConvertDiceRolls | 27,556 | 4,112 | **6.7x faster** | 69,732 | 3,261 | 763 | 65 |
| ExtractModifier | 585.3 | 79.02 | **7.4x faster** | 1,480 | 0 | 16 | 0 |
| IsD20Roll | 2,047 | 67.71 | **30.2x faster** | 5,338 | 0 | 69 | 0 |
| CalculateAverage | 3,048 | 517.9 | **5.9x faster** | 8,088 | 352 | 92 | 11 |
| ParseDiceNotation | 2,796 | 177.6 | **15.7x faster** | 7,956 | 138 | 91 | 5 |
| ConvertDiceRollsVaried | 21,134 | 2,993 | **7.1x faster** | 52,797 | 2,150 | 566 | 45 |

## Key Findings

### Speed Improvements
- **Best improvement**: `IsD20Roll()` is now **30.2x faster** (from 2,047 ns/op to 67.71 ns/op)
- **Overall improvement**: All functions show **5.9x to 30.2x speedup**
- **Real-world scenario**: `ConvertDiceRolls()` improved **6.7x** (27.5 µs → 4.1 µs per call)

### Memory Improvements
- **Zero allocations**: `ExtractModifier()` and `IsD20Roll()` now have **0 allocations**
- **Memory reduction**: Up to **95%+ reduction** in memory allocations
- **Example**: `ConvertDiceRolls()` reduced from 69,732 B/op to 3,261 B/op (**95.3% reduction**)

### Allocation Reduction
- **Dramatic reduction**: Functions reduced allocations from 16-92 down to 0-11
- **ExtractModifier**: 16 allocs → 0 allocs
- **IsD20Roll**: 69 allocs → 0 allocs
- **ConvertDiceRolls**: 763 allocs → 65 allocs (**91.5% reduction**)

## Implementation Details

### Changed Regexes
The following regular expressions were moved from function-level to package-level variables:

```go
var (
    diceNotationRegex = regexp.MustCompile(`^(\d*)d(\d+)([+-]\d+)?$`)
    rollPatternRegex  = regexp.MustCompile(`(to hit|damage|healing|save):\s*(\d*d\d+[+-]?\d*)`)
    modifierRegex     = regexp.MustCompile(`[+-]\d+`)
    d20RollRegex      = regexp.MustCompile(`^\d*d20([+-]\d+)?$`)
    averageRegex      = regexp.MustCompile(`^(\d+)d(\d+)([+-]\d+)?$`)
)
```

### Functions Optimized
1. `ParseDiceNotation()` - Uses `diceNotationRegex`
2. `ConvertDiceRolls()` - Uses `rollPatternRegex`
3. `extractModifier()` - Uses `modifierRegex`
4. `isD20Roll()` - Uses `d20RollRegex`
5. `calculateAverage()` - Uses `averageRegex`

## Impact on Real-World Usage

For a typical character sheet with 10 abilities, each containing 2-3 dice rolls:
- **Before**: ~825 µs total parsing time
- **After**: ~123 µs total parsing time
- **Savings**: ~700 µs per character sheet (**6.7x faster**)

For a batch operation processing 1,000 character sheets:
- **Before**: ~825 ms
- **After**: ~123 ms
- **Savings**: ~702 ms (**85% faster**)

## Benchmark Environment

- **CPU**: Apple M3 Pro
- **OS**: darwin/arm64
- **Go Version**: go1.23 or later
- **Benchmark Time**: 3s per benchmark
- **Date**: 2026-02-18

## Conclusion

Moving regex compilation to package-level variables eliminates the overhead of repeatedly compiling the same patterns. This is especially impactful for:

1. Functions called in loops (like `ConvertDiceRolls()` for each ability)
2. Simple validation functions (like `isD20Roll()` and `extractModifier()`)
3. Functions called multiple times per operation (like `calculateAverage()`)

The optimization is a clear win with:
- **Zero code complexity increase**
- **No API changes**
- **Significant performance gains**
- **Dramatic memory reduction**
- **All tests passing**

This demonstrates the importance of profiling and benchmarking, as regex compilation overhead can be a hidden performance bottleneck in frequently-called functions.
