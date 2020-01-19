const request = require('request-promise-native');

async function list() {
  try {
    const res = await request.get('http://localhost:8080/list', { json: true });
    let streams = res.streams || [];
    if (!streams.length) {
      console.log(`No streams found`);
    }
    else {
      console.log(`Found the following streams:`)
      console.log(streams.join('\n'));
    }
  } catch (err) {
    console.error("Error: ", err);
  }
}

module.exports = list;