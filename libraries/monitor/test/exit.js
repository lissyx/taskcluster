const MonitorManager = require('../');

const options = {
  shouldError: false,
  shouldUnhandle: false,
  patchGlobal: false,
  bailOnUnhandledRejection: false,
};

process.argv.slice(2).forEach(arg => {
  options[arg.slice(2)] = true;
});

const manager = new MonitorManager({
  serviceName: 'foo-testing',
});
manager.setup(options);

if (options.shouldError) {
  throw new Error('hello there');
}
if (options.shouldUnhandle) {
  Promise.reject('whaaa');
}
