const program = require('commander');
const status = require('./lib/status');
const list = require('./lib/list');
const stream = require('./lib/stream');

module.exports = () => {
  program
    .command('status')
    .description('Check if the daemon is runnning')
    .action(() => {
      status();
    });

  program
    .command('list')
    .description('Print the available container log records')
    .action(() => {
      list();
    });

  program
    .command('stream <streamId>')
    .description('Streams logs for the specified streamId')
    .action((streamId) => {
      stream(streamId);
    });

  program.parse(process.argv);
}