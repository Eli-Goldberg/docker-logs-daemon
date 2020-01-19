const Logger = require('./logger');

const APP_NAME = process.env.APP_NAME || "Example App";
const log = Logger(APP_NAME);

const LOG_WRITE_INTERVAL = 500 // 0.5 second
// const LOG_WRITE_PERIOD = 5 * 60 * 1000; // 1 minute



const randomLogActions = [
  () => log.info('Some Info Msg'),
  () => log.warn('Some Warn Msg'),
  () => log.error('Some Error Msg')
];

const getRandomInt = (max) => Math.floor(Math.random() * Math.floor(max));
const interval = setInterval(() => {
  const logAction = randomLogActions[getRandomInt(randomLogActions.length)];
  logAction();
}, LOG_WRITE_INTERVAL);

// setTimeout(() => {
//   clearImmediate(interval);
//   console.log(`Done, exiting...`)
// }, LOG_WRITE_PERIOD);

