// Declaration file for Go's math package
declare module "go:math" {
    // Constants
    const Pi: float
    const E: float
    const Phi: float
    const Sqrt2: float
    const SqrtE: float
    const SqrtPi: float
    const Ln2: float
    const Ln10: float
    const MaxFloat64: float
    const MaxInt: int
    const MinInt: int

    // Basic functions
    function Sqrt(x: float): float
    function Abs(x: float): float
    function Floor(x: float): float
    function Ceil(x: float): float
    function Round(x: float): float
    function Trunc(x: float): float

    // Trigonometric functions
    function Sin(x: float): float
    function Cos(x: float): float
    function Tan(x: float): float
    function Asin(x: float): float
    function Acos(x: float): float
    function Atan(x: float): float
    function Atan2(y: float, x: float): float
    function Sinh(x: float): float
    function Cosh(x: float): float
    function Tanh(x: float): float

    // Exponential and logarithmic functions
    function Log(x: float): float
    function Log10(x: float): float
    function Log2(x: float): float
    function Exp(x: float): float
    function Exp2(x: float): float
    function Pow(x: float, y: float): float
    function Pow10(n: int): float

    // Comparison functions
    function Max(x: float, y: float): float
    function Min(x: float, y: float): float
    function Mod(x: float, y: float): float
    function Remainder(x: float, y: float): float

    // Special functions
    function IsNaN(f: float): boolean
    function IsInf(f: float, sign: int): boolean
    function NaN(): float
    function Inf(sign: int): float
    function Signbit(x: float): boolean
    function Copysign(f: float, sign: float): float
}
