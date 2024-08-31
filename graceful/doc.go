// Package graceful implements functions for handling graceful shutdown.
//
// The [NewGracefulShutdown] function creates a graceful shutdown handler.
//
// Life cycle: graceful shutdown handler
//   - calls the function to start each [GracefulServer]
//   - if there is an error starting someone [GracefulServer], the graceful shutdown handler will initiate its shutdown
//   - will wait for an interrupt signal
//   - after an interrupt signal, calls the function to stop each [GracefulServer]
//   - will wait for all [GracefulServer] to stop (respecting timeout if configured)
//
// Life cycle: each [GracefulServer]
//   - will wait for the graceful shutdown handler to signal its shutdown
//   - call the function to stop the [GracefulServer]
//   - will wait for the [GracefulServer] to shut down
//   - if timeout is set, it will call the function to forcibly stop [GracefulServer] (it will not wait for this process)
//   - informs the graceful shutdown handler that the [GracefulServer] has been shut down
package graceful
