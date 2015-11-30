/*

testd that makes easier to write integration tests that depends on runnings daemons/services.

Example:

  // Create a Testd instance, informing where you want the logs to be saved,
  // the command to execute and its arguments.
  //
  daemon, err := testd.New("./tests/debug/daemon.log", "daemon", "arg1", "arg2")

  // Do your testing, that depends on the daemon being running.
  // When you are done testing, just call the Stop method.
  err := daemon.Stop()

  // Remember to do proper error handling, both on New and Stop.
*/
package testd
