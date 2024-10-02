// The graceful package provides functionality to handle graceful startup and shutdown processes for servers.
// It ensures that servers start and stop in a controlled manner, handling tasks such as waiting for active requests to
// complete before shutting down, listening for interrupt signals, and managing timeout scenarios.
// The package is designed to simplify the management of multiple servers, allowing them to be started and stopped
// gracefully while minimizing disruptions.
//
// # Key Features
//   - Supports the graceful shutdown of servers by waiting for ongoing requests to finish.
//   - Manages multiple servers, allowing them to be started and stopped.
//   - Interrupt signal handling, listens for system signals (SIGINT, SIGTERM) to initiate the shutdown process.
//   - Configurable timeout support for shutdown operations, with forced stop functionality if timeout is exceeded.
//   - Error management, automatically triggers shutdown procedures if any server encounters startup errors, maintaining application stability.
//   - Once-executed guarantee, The [GracefulShutdown] interface is executed exactly once to avoid conflicting shutdown actions.
//
// # Life Cycle
//
// 1. Initialization Phase
//   - A [GracefulShutdown] handler is created using [NewGracefulShutdown] and requires registering one or more [GracefulServer] instances.
//   - For each server, a [GracefulServer] is created using [NewGracefulServer].
//   - A [GracefulServer] specifically for an [http.Server] is created using [NewGracefulServerHttp].
//
// 2. Startup Phase
//   - When the Run method on [GracefulShutdown] is called, starts all registered servers by calling each server's Start method.
//   - Each server begins processing requests as per its defined behavior.
//   - The application now enters its normal operational state where servers are running and handling requests.
//   - Error handling if any [Graceful Server] fails to start, it will initiate the shutdown process.
//
// 3. Waiting Phase
//   - During this phase, the [GracefulShutdown] handler waits for an interrupt signal (a SIGINT or SIGTERM) or a
//     cancellation/timeout signal from the provided context.
//   - The handler remains idle, letting the servers run until such a signal is received.
//
// 4. Shutdown Initiation
//   - When an interrupt signal or a context cancellation occurs, the shutdown process begins.
//   - The [GracefulShutdown] handler invokes the registered [WithNotifyShutdown] function to notify that the shutdown process has begun.
//   - The [GracefulShutdown] handler iterates through the registered [GracefulServer] instances and calls their Stop method,
//     allowing the server to finish processing ongoing requests before shutting down.
//   - The [GracefulShutdown] handler will continue waiting for servers to complete their shutdown within the allotted time (if a timeout was set).
//
// 5. Timeout Handling (Optional)
//   - If a timeout is defined and any server has not yet completed its graceful shutdown within that time, the handler calls ForceStop on that server.
//   - Force stop should immediately shut down the server, regardless of any ongoing requests.
//
// 6. Completion Phase
//   - After all servers have either shut down gracefully or been forcefully stopped, the [GracefulShutdown] handler completes its lifecycle.
//   - The application is now fully stopped, and the lifecycle ends.
//
// 7. Post-Shutdown Actions
//   - Perform any necessary cleanup activities before fully exiting the application.
package graceful
