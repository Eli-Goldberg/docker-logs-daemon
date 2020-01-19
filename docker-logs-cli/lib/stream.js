const WebSocket = require('ws');

async function stream(streamId) {
  try {
    console.log(`Connecting to daemon log stream for stream id: ${streamId}`);
    const ws = new WebSocket('ws://localhost:8080/ws');
    ws.on('open', function open() {
      console.log('connected');
      ws.send(JSON.stringify({ StreamID: streamId }));
    });

    ws.on('close', function close() {
      console.log('disconnected');
    });

    ws.on('message', function incoming(message) {
      console.log(message);
    });
    ws.on('error', function () {
      console.error('Error: Lost connection')
    });

  } catch (err) {
    console.error("Error: ", err);
  }
}

module.exports = stream;