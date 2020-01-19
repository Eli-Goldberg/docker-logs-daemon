const bunyan = require('bunyan');
const safeCycles = bunyan.safeCycles;

class BunyanSpecificLevelStream {
  constructor(levels, stream) {
    this.stream = stream;
    const _levels = {};
    levels.forEach(function (lvl) {
        _levels[bunyan.resolveLevel(lvl)] = true;
    });
    this.levels = _levels;
  }
  write(rec) {
    if (this.levels[rec.level] !== undefined) {
      const str = JSON.stringify(rec, safeCycles()) + '\n';
      this.stream.write(str);
    }
  }
}

module.exports = BunyanSpecificLevelStream;