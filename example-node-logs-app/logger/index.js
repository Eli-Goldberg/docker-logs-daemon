const bunyan = require('bunyan');
const BunyanSpecificLevelStream = require('./BunyanSpecificLevelStream');

const createLogger = (logName) => bunyan.createLogger({
  name: logName,
  streams: [
    {
      level: 'info',
      type: 'raw',
      // log only INFO, WARN and above to stderr
      stream: new BunyanSpecificLevelStream(['info', 'warn'], process.stdout)
    },
    {
      level: 'error',
      // log ERROR and above to stderr
      stream: process.stderr
    }]
});


module.exports = createLogger;