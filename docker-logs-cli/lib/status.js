const request = require('request-promise-native');

async function status() {
  let ok = false;
  try {
    const res = await request.get('http://localhost:8080/status', { json: true });
    ok = (res.ok === true);
  } catch (err) {}
  
  if (ok) {
    console.log("Daemon is up and running")
  } else {
    console.error("Daemon is down");
  }
}

module.exports = status;