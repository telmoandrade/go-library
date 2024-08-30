// Package graceful implements functions to manipulate graceful shutdown.
//
// The [NewGracefulShutdown] function creates a graceful shutdown manager.
//
// graceful shutdown [Run]
// - starts each registered server
// - will wait for an interruption signal
// - stops running each server
// - will wait for all servers to finish
//
// if there is an error when starting a server, the manager will be informed that it must be shut down
//
// each server
// - will wait for the manager to signal its shutdown
// - stop the server
// - will wait for the server to shut down
// - if the timeout is configured, it will forcibly stop the server (it will not wait for this process)
// - informs the manager that the server has been shut down
package graceful
