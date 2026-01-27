# goTS Built-in Functions TODO

This document tracks the built-in functions and methods that goTS needs to implement to achieve better JavaScript/TypeScript compatibility. The list is derived from the ECMAScript test262 conformance test suite.

## Currently Implemented

### Global Functions
- `println`, `print` - Console output
- `len` - Get length of arrays/strings
- `push`, `pop` - Array mutation (functional style)
- `typeof` - Type checking
- `tostring`, `toint`, `parseInt`, `tofloat` - Type conversion
- `sqrt`, `floor`, `ceil`, `abs` - Basic math (legacy global functions)

### Math Object (NEW)
**Constants:**
- `Math.PI`, `Math.E`

**Methods:**
- Rounding: `round`, `floor`, `ceil`, `trunc`
- Power/Roots: `pow`, `sqrt`, `cbrt`, `exp`
- Logarithms: `log`, `log10`, `log2`
- Absolute/Sign: `abs`, `sign`
- Min/Max: `min`, `max` (variadic)
- Trigonometry: `sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2`
- Random: `random`

### Array Methods
- Mutation: `push`, `pop`, `shift`, `unshift`, `splice`, `reverse`, `sort`
- Access: `slice`, `concat`, `indexOf`, `includes`, `join`
- Iteration: `map`, `filter`, `reduce`, `forEach`, `find`, `findIndex`, `some`, `every`
- Properties: `length`

### String Methods
- `split`, `replace`, `trim`
- `startsWith`, `endsWith`, `includes`
- `toLowerCase`, `toUpperCase`
- `substring`, `charAt`, `indexOf`

### Map Methods
- `get`, `set`, `has`, `delete`, `clear`
- `keys`, `values`, `entries`, `forEach`
- Properties: `size`

### Set Methods
- `add`, `has`, `delete`, `clear`, `values`, `forEach`
- Properties: `size`

### RegExp Methods
- `test`, `exec`

### Number Methods
- `toString`

---

## Missing Built-ins (Priority Order)

### Priority 1: Essential Core Functions (PARTIALLY DONE)

#### Math Object (Static Methods) - IMPLEMENTED
~~Math object methods~~ - Now implemented via built-in object registry pattern.

#### Number Static Methods
```typescript
Number.isFinite(x)     // Check if finite number
Number.isNaN(x)        // Check if NaN
Number.isInteger(x)    // Check if integer
Number.isSafeInteger(x)// Check if safe integer
Number.parseFloat(s)   // Parse float from string
Number.parseInt(s, r)  // Parse int from string with radix
// Constants
Number.MAX_VALUE
Number.MIN_VALUE
Number.POSITIVE_INFINITY
Number.NEGATIVE_INFINITY
Number.NaN
```

#### Number Prototype Methods
```typescript
num.toFixed(digits)    // Format with fixed decimals
num.toPrecision(p)     // Format with precision
num.toExponential(d)   // Exponential notation
```

#### Global Functions
```typescript
isNaN(x)               // Global NaN check
isFinite(x)            // Global finite check
parseFloat(s)          // Global float parser
parseInt(s, radix)     // Global int parser (partial - have toint)
encodeURI(s)           // Encode URI
decodeURI(s)           // Decode URI
encodeURIComponent(s)  // Encode URI component
decodeURIComponent(s)  // Decode URI component
```

### Priority 2: String Methods

```typescript
// Access
str.charCodeAt(i)      // Unicode code point at index
str.codePointAt(i)     // Full Unicode code point
str.at(i)              // Element at index (negative allowed)

// Search
str.search(regexp)     // Search with regex, return index
str.match(regexp)      // Match against regex
str.matchAll(regexp)   // All matches as iterator

// Transform
str.slice(start, end)  // Extract section (negative indices)
str.repeat(count)      // Repeat string
str.padStart(len, s)   // Pad at start
str.padEnd(len, s)     // Pad at end
str.replaceAll(s, r)   // Replace all occurrences
str.normalize(form)    // Unicode normalization

// Trim variants
str.trimStart()        // Trim leading whitespace
str.trimEnd()          // Trim trailing whitespace

// Case
str.toLocaleLowerCase()
str.toLocaleUpperCase()

// Comparison
str.localeCompare(s)   // Locale-aware comparison
```

### Priority 3: Array Methods

```typescript
// Static methods
Array.isArray(x)       // Check if array
Array.from(iter)       // Create from iterable
Array.of(...items)     // Create from arguments

// Access
arr.at(i)              // Element at index (negative allowed)
arr.flat(depth)        // Flatten nested arrays
arr.flatMap(fn)        // Map then flatten

// Search
arr.lastIndexOf(v)     // Last index of value
arr.findLast(fn)       // Find from end
arr.findLastIndex(fn)  // Find index from end

// Immutable variants (ES2023)
arr.toReversed()       // Return reversed copy
arr.toSorted(fn)       // Return sorted copy
arr.toSpliced(...)     // Return spliced copy
arr.with(i, v)         // Return copy with replaced element

// Copy
arr.copyWithin(t,s,e)  // Copy within array
arr.fill(v, s, e)      // Fill with value
```

### Priority 4: Object Static Methods

```typescript
Object.keys(obj)           // Get own property names
Object.values(obj)         // Get own property values
Object.entries(obj)        // Get [key, value] pairs
Object.fromEntries(iter)   // Create from entries
Object.assign(target, ...sources)  // Copy properties
Object.freeze(obj)         // Make immutable
Object.seal(obj)           // Prevent new properties
Object.isFrozen(obj)       // Check if frozen
Object.isSealed(obj)       // Check if sealed
Object.hasOwn(obj, prop)   // Check own property
Object.create(proto)       // Create with prototype
Object.getPrototypeOf(obj) // Get prototype
Object.setPrototypeOf(o,p) // Set prototype
```

### Priority 5: JSON

```typescript
JSON.parse(text)           // Parse JSON string
JSON.stringify(value)      // Convert to JSON string
JSON.stringify(v, replacer, space)  // With formatting
```

### Priority 6: Date Object

```typescript
// Constructor
new Date()
new Date(ms)
new Date(dateString)
new Date(y, m, d, h, min, s, ms)

// Static
Date.now()             // Current timestamp
Date.parse(s)          // Parse date string
Date.UTC(...)          // UTC timestamp

// Getters
date.getTime()
date.getFullYear()
date.getMonth()
date.getDate()
date.getDay()
date.getHours()
date.getMinutes()
date.getSeconds()
date.getMilliseconds()
date.getTimezoneOffset()
// UTC variants: getUTCFullYear, getUTCMonth, etc.

// Setters
date.setTime(ms)
date.setFullYear(y)
date.setMonth(m)
date.setDate(d)
date.setHours(h)
date.setMinutes(m)
date.setSeconds(s)
date.setMilliseconds(ms)
// UTC variants: setUTCFullYear, setUTCMonth, etc.

// Formatting
date.toString()
date.toDateString()
date.toTimeString()
date.toISOString()
date.toJSON()
date.toLocaleString()
date.toLocaleDateString()
date.toLocaleTimeString()
```

### Priority 7: Promise (Async Support)

```typescript
// Constructor
new Promise((resolve, reject) => { })

// Static
Promise.resolve(value)
Promise.reject(reason)
Promise.all(iterable)
Promise.allSettled(iterable)
Promise.race(iterable)
Promise.any(iterable)

// Prototype
promise.then(onFulfilled, onRejected)
promise.catch(onRejected)
promise.finally(onFinally)
```

### Priority 8: Additional Collections

#### WeakMap
```typescript
weakMap.get(key)
weakMap.set(key, value)
weakMap.has(key)
weakMap.delete(key)
```

#### WeakSet
```typescript
weakSet.add(value)
weakSet.has(value)
weakSet.delete(value)
```

### Priority 9: Typed Arrays

```typescript
// Constructors
new Uint8Array(length)
new Int8Array(length)
new Uint16Array(length)
new Int16Array(length)
new Uint32Array(length)
new Int32Array(length)
new Float32Array(length)
new Float64Array(length)
// BigInt variants
new BigInt64Array(length)
new BigUint64Array(length)

// ArrayBuffer
new ArrayBuffer(byteLength)
arrayBuffer.slice(begin, end)
ArrayBuffer.isView(arg)

// DataView
new DataView(buffer)
dataView.getInt8(offset)
dataView.setInt8(offset, value)
// ... other typed access methods
```

### Priority 10: Symbol

```typescript
Symbol(description)
Symbol.for(key)
Symbol.keyFor(sym)

// Well-known symbols
Symbol.iterator
Symbol.asyncIterator
Symbol.toStringTag
Symbol.hasInstance
// etc.
```

### Priority 11: Reflect & Proxy

```typescript
// Reflect
Reflect.get(target, prop)
Reflect.set(target, prop, value)
Reflect.has(target, prop)
Reflect.deleteProperty(target, prop)
Reflect.ownKeys(target)
Reflect.apply(fn, thisArg, args)
Reflect.construct(target, args)

// Proxy
new Proxy(target, handler)
Proxy.revocable(target, handler)
```

### Priority 12: Error Types

```typescript
new Error(message)
new TypeError(message)
new RangeError(message)
new SyntaxError(message)
new ReferenceError(message)
new EvalError(message)
new URIError(message)
new AggregateError(errors, message)

// Properties
error.name
error.message
error.stack
error.cause
```

---

## Implementation Notes

### Go Mappings

| JavaScript | Go Equivalent |
|------------|---------------|
| `Math.random()` | `rand.Float64()` |
| `Math.round(x)` | `math.Round(x)` |
| `Math.pow(x, y)` | `math.Pow(x, y)` |
| `Math.min/max` | Custom variadic func |
| `Number.isNaN(x)` | `math.IsNaN(x)` |
| `Number.isFinite(x)` | `!math.IsInf(x, 0) && !math.IsNaN(x)` |
| `JSON.parse()` | `json.Unmarshal()` |
| `JSON.stringify()` | `json.Marshal()` |
| `Date` | `time.Time` |
| `Promise` | Goroutines + channels |
| `Map` | `map[K]V` (done) |
| `Set` | `map[T]struct{}` (done) |
| `WeakMap/WeakSet` | Not directly possible (GC) |
| `Symbol` | String constants or iota |
| `Proxy/Reflect` | Interfaces + reflection |
| `TypedArray` | `[]byte`, `[]int32`, etc. |

### Challenges

1. **WeakMap/WeakSet**: Go's GC doesn't support weak references natively. May need to use `sync.Map` with explicit cleanup or accept that items won't be garbage collected.

2. **Symbols**: Go doesn't have symbols. Can simulate with unique strings or use constants.

3. **Proxy**: Requires Go's `reflect` package extensively. Performance may differ significantly.

4. **async/await**: Requires goroutines and channels. Current Promise support uses generic type `GTS_Promise[T]`.

5. **Date**: Go's `time.Time` is different from JavaScript's Date. Month indexing differs (Go: 1-12, JS: 0-11).

---

## Progress Tracking

- [x] Math object (20+ methods) - DONE
- [x] Number static/prototype methods (10+ methods) - DONE
- [ ] String additional methods (15+ methods)
- [ ] Array additional methods (10+ methods)
- [ ] Object static methods (15+ methods)
- [x] JSON parse/stringify - DONE
- [ ] Date object
- [ ] Promise enhancements
- [ ] WeakMap/WeakSet
- [ ] TypedArrays
- [ ] Symbol basics
- [ ] Reflect/Proxy
- [ ] Error types

---

## Test262 Reference

The complete ECMAScript test suite is in `scaffold/test262/`. Key directories:
- `test/built-ins/` - Built-in object tests
- `test/language/` - Language feature tests
- `harness/` - Test utilities

Use these tests as reference for expected behavior when implementing built-ins.
