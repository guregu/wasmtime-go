package wasmtime

// #cgo !windows LDFLAGS:-lwasmtime
// #cgo windows LDFLAGS:-lwasmtime.dll
import "C"

// # What's up with `ptr()` methods?
//
// We use `runtime.SetFinalizer` to free all objects we allocate from C. This
// is intended to make usage of the API much simpler since you don't have to
// close/free anything. The tricky part here though is laid out in
// `runtime.SetFinalizer`'s documentation which is that if you read a
// non-gc-value (like a C pointer) from a GC object then after the value is
// read the GC value might get garbage collected. This is quite bad for us
// because the garbage collection will free the C pointer, making the C pointer
// actually invalid.
//
// The solution is to add `runtime.KeepAlive` calls after C function calls to
// ensure that the GC object lives at least as long as the C function call
// itself. This is naturally quite error-prone, so the goal here with `ptr()`
// methods is to make us a bit more resilient to these sorts of errors and
// expose segfaults during development.
//
// Each `ptr()` method has the basic structure of doing these steps:
//
// 1. First it reads the pointer value from the GC object
// 2. Next it conditionally calls `runtime.GC()`, depending on build flags
// 3. Finally it returns the original pointer value
//
// The goal here is to as aggressively as we can collect GC objects when
// testing and trigger finalizers as frequently as we can. This naturally
// slows things down quite a bit, so conditional compilation (with the `debug`
// tag) is used to enable this. Our CI runs tests with `-tag debug` to make
// sure this is at least run somewhere.
//
// If anyone else has a better idea of what to handle all this it would be very
// much appreciated :)
