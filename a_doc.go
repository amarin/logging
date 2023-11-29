/*
Package logging provides standard interface for Logger and common Config structure to use with logging.
It contains no implementation of logging itself.
Logging functions are enabled with simple wrapper over actual logging packages, f.e. Zap or Logrus.

Concrete logging backend is up to developer.
Chosen backend should be wrapped into Backend implementation if not implements it directly.

Some standard logging keys also provided to follow single logging style everywhere.
*/
package logging
