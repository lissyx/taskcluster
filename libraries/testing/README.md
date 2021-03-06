# Testing Library

Support for testing Taskcluster components.

This module contains a number of utilities that facilitate testing Taskcluster
components.  It is typically installed as a devDependency, so it is not used in
production code.

See the source for detailed documentation.

Sticky Loader
-------------

A sticky loader is a thin wrapper around `taskcluster-lib-loader` to support
dependency injection. It "remembers" each value it has returned and will return
it again on the next call; it can also have a dependency injected.  Use it like
this in `helper.js`:

```javascript
const {stickyLoader} = require('taskcluster-lib-testing');
const load = require('../src/server');

exports.load = stickyLoader(load);
exports.load.inject('profile', 'test');
exports.load.inject('process', 'test');
```

The `load.inject(component, value)` method sets a loader overwrite without
attempting to load it. There is a corresponding `load.remove(component)` to
remove a component.

In test scripts:

```javascript
const {load} = require('./helper');

suite('SomeTable', function() {
  suiteSetup(async function() {
    load.save(); // save the state of the loader to restore in tearDown
    await load('cfg'); // load the cfg so we can edit it
    load.cfg('azure.accountName', 'inMemory'); // edit the cfg in-place
    const SomeTable = await load('SomeTable');
    await SomeTable.ensureTable({ /* ... */ });
  });

  suiteTeardown(function() {
    load.restore(); // restore the state of the loader
  });

  test(async function() {
    const component = await load('some-component');
    // some-component will be loaded with the same cfg and with
    // the same instance of SomeTable that we set up above
  });
});
```


The `load.save()` and `load.restore()` methods push and pull loader states in a
stack, and are best used in setup/teardown methods to ensure that one suite
does not "pollute" the loader state for the next.

The `load.cfg(path, value)` method edits the `cfg` component in place, using a
dotted path to specify the config value. The `save` and `restore` methods are
careful to deep-copy `cfg` so that these in-place modifications affect only
the current loader state.

If `cfg` is not loaded, the `load.cfg()` method will not work, so generally (as
in the example above) a bare `load('cfg')` is used to ensure its presence.

Secrets
-------

This class handles getting secrets for tests, and easily determining what
secrets are available.  It integrates with `taskcluster-lib-config`.  Set it up by
in `test/helper.js`:

```javascript
const {Secrets} = require('taskcluster-lib-testing');

exports.secrets = new Secrets({
  secretName: [
    'project/taskcluster/testing/taskcluster-foo',
    'project/taskcluster/testing/taskcluster-foo/master-only',
  ],
  // provide a stickyLoader instance for use in mockSuite
  load,
  secrets: {
   pulse: [
     // env - the environment variable by which this secret is set in the config (if any)
     // cfg - dotted path to the config value containing this secret (if any)
     // name - name for the secret (used for programmatic access in tests; defaults to env)
     // mock - value to provide if secret is not set (for mock runs only)
     {env: 'PULSE_USERNAME', cfg: 'pulse.username', name: 'username', mock: 'dummy'},
     {env: 'PULSE_PASSWORD', cfg: 'pulse.password', name: 'password'},
   ],
   aws: [
     {env: 'AWS_ACCESS_KEY_ID', cfg: 'aws.accessKeyId'},
     {env: 'AWS_SECRET_ACCESS_KEY', cfg: 'aws.secretAccessKey'},
   ],
  },
});
```

If a secret is defined in the loaded configuration, that value will be used even if the `env` key is also set.
Secrets should not have any value set in `config.yml` (although `!env` is OK), or this class will not function properly.
If the system you are testing does not use `taskcluster-lib-config`, simply do not specify the `cfg` properties to the constructor.

You can then call `await secrets.setup()`  to set up the secrets (reading from `cfg` if necessary).
This *must* be called during Mocha's runtime, so either in a setup function or a test.
It short-circuits multiple calls, so it's safe to call it all over the place.
In fact, `mockSuite` (below) will call it for you.

In CI (when `$TASK_ID` is set), the `setup` method will attempt to fetch the secrets named in `secretName` from the secrets service.
It expects the fetch value to be a map from environment variable name to value.
If a fetch fails, it is considered equivalent to fetching an empty map.
This allows, for example, secrets that can only be fetched on pushes to the master branch, and not pull requests.

The secrets object has a few useful methods, all of which can only be called *after* `setup`, and thus only in a setup function or a test:

* `secrets.have(name)` -- true if the given secret is available
* `secrets.get(name)` -- returns an object containing the secret values by name, or throws an error if not avaialble

## mockSuite

The `secrets.mockSuite` function abstracts away the most common case: running the same tests in a mock and real environment, skipping the real tests if secrets are not available.
It is called as `secrets.mockSuite(title, [secrets], async function(mock, skipping) { .. })` in the same location you might call Mocha's `suite(..)`.
The `secrets` is an array of secret names required to run this suite in a real environment.
The given function should define the suite, and can include `setup`, `suiteSetup`, and so on.
The `mock` parameter is true for the mock version, and false for the real version.
If `$NO_TEST_SKIP` is set, `mockSuite` will throw an error when secrets are not available.

Note that Mocha continues to run `setupSuite` and `teardownSuite` functions even after a suite has been skipped.
Mocha does not provide any way to determine if a suite has been skipped.
Use `skipping()` to determine if the suite is currently skipping, and avoid doing initialization that will fail.

Note, too, that all modern versions of Mocha have [a bug](https://github.com/mochajs/mocha/issues/2819) causing nested suites to run anyway, even when the parent suite is skipped.
A quick (but unfortunate) way to work around this bug is

```javascript
suite('mySuite', function() {
  suiteSetup(function() {
    if (skipping()) {
      this.skip();
    }
  });
});
```

### Usage

```javascript
// helper.js
const {Secrets, stickyLoader} = require('taskcluster-lib-testing');

const load = stickyLoader(require('../src/main'));
const secrets = new Secrets({
  secretName: 'project/taskcluster/testing/taskcluster-ping',
  secrets: {
    pingdom: [
      {name: 'apiKey', env: 'PINGDOM_API_KEY', cfg: 'app.pingdom.apiKey'},
    ],
    taskcluster: [
      {name: 'clientId', env: 'TASKCLUSTER_CLIENT_ID', cfg: 'taskcluster.credentials.clientId'},
      {name: 'accessToken', env: 'TASKCLUSTER_ACCESS_TOKEN', cfg: 'taskcluster.credentials.accessToken'},
    ],
  },
  load,
});

exports.secrets = secrets;
exports.load = load;
```

```javascript
// some_test.js
const {secrets, load} = require('./helper');

// for testing by passing secrets to the subject..
secrets.mockSuite('pingdom updates', ['pingdom'], function(mock, skipping) {
  let pingdomUpdater, pingdomComponent;

  suiteSetup(async function() {
    // use secrets.get(..) in the real case
    pingdomUpdater = new PingdomUpdater({apiKey: mock ? 'pretendKey' : secrets.get('pingdom').apiKey});
    if (mock) {
      nock('https://pingdom.com:443', ..); // mock out Pingdom API
    }
  });

  suiteTeardown(function() {
    if (mock) {
      nock.clearAll();
    }
  });

  test('updates once', function() { .. });
});

// for testing a loader component..
secrets.mockSuite('Floobits', ['taskcluster'], function(mock) {
  let Floobits;
  suiteSetup(async function() {
    if (mock) {
      // set the special accountName that will cause azure-entities to use its fake version;
      // otherwise, the loader component will use the taskcluster secrets to get access
      // to the a Azure table
      helper.load.cfg('azure.accountName', 'inMemory');
    }
    
    if (!skipping()) {
      Floobits = await helper.load('Floobits');
      await Floobits.ensureTable();
    }
  });

  test('create', async function() {
    await Floobits.create(..);
    // ..
  });
});
```

The test output for the first suite will contain something like

```
  pingdom updates (mock)
    ✓ updates once
  pingdom updates (real)
    - updates once
```


PulseTestReceiver
-----------------

A utility for tests written in mocha, that makes it very easy to wait for a
specific pulse message.  This uses real pulse messages, so pulse credentials
will be required.

**Example:**
```js
suite("MyTests", function() {
  let credentials = {
    username:     '...',  // Pulse username
    password:     '...'   // Pulse password
  };
  let receiver = new testing.PulseTestReceiver(credentials, mocha)

  test("create task message arrives", async function() {
    var taskId = slugid.v4();

    // Start listening for a message with the above taskId, giving
    // it a local name (here, `my-create-task-message`)
    await receiver.listenFor(
      'my-create-task-message',
      queueEvents.taskCreated({taskId: taskId})
    );

    // We are now listen for a message with the taskId
    // So let's create a task with it
    await queue.createTask(taskId, {...});

    // Now we wait for the message to arrive
    let message = await receiver.waitFor('my-create-task-message');
  });
});
```

The `receiver` object will setup an PulseConnection before all tests and close
the PulseConnection after all tests. This should make tests run faster.  All
internal state, ie. the names given to `listenFor` and `waitFor` will be reset
between all tests.

schemas
-------

Test schemas with a positive and negative test cases.

The method should be called within a `suite`, as it will call the mocha `test`
function to define a test for each schema case.

 * `schemasetOptions` - {}  // options to pass to the [taskcluster-lib-validate](https://github.com/taskcluster/taskcluster-lib-validate) constructor
 * `cases` - array of test cases
 * `basePath` -  base path for relative pathnames in test cases (default `path.join(__dirname, 'validate')`)

Each test case looks like this:

```js
{
  schema:   'https://tc-tests.localhost/svcname/v7/frobnicate-foo.json', // JSON schema identifier to test against
  path:     'test-file.json',             // Path to test file (relative to basePath)
  success:  true || false                 // true if validation should succeed; false if it should fail
}
```

fakeauth
--------

A fake for the auth service to support testing APIs without requiring
production credentials, using Nock.

This object intercepts requests to the auth service's `authenticateHawk` method
and returns a response based on the given `clients`, instead. Note that
accessTokens are not checked -- the fake simply controls access based on
clientId or the scopes in a temporary credential or supplied with
authorizedScopes.

To start the mock, call `testing.fakeauth.start(clients, {rootUrl})` in your suite's
`setup` method. The first argument has the form

```js
{
 "clientId1": ["scope1", "scope2"],
 "clientId2": ["scope1", "scope3"],
}
```

The auth service on the cluster identified by `rootUrl` will be faked. When
used to test an API in a microservice, this is same as the root URL for the
fake web server -- `http://localhost:1234` or something of that sort.

Call `testing.fakeauth.stop()` in your test suite's `teardown` method to stop the HTTP interceptor.

withEntity
----------

This function is intended for use with `mockSuite` and the usual Azure configuration for services.
Given a list of entity classes, it either loads an in-memory version (for `mock = true`) or uses a unique table name.
It assigns the set-up class instance to `helper.<name>` for easy access from tests.

This is typically called from a service's `helper.js` like this:

```js
exports.withEntities = async (mock, skipping) => {
  withEntity(mock, skipping, exports, 'Thing', data.Thing, {});
});
```

where `exports` is the `helper` module and `Thing` is the table being loaded, usually defined in `src/data.js`.
Ideally the class name and the loader component match, but in cases where the same class is used for multiple tables they can differ.
This call can be repeated for additional tables.

The last argument contains options.
For cases where tests assume that the table state persists between test cases, pass `{orderedTests: true}` as an option in the final argument.
For cases where the table-cleanup operation is more complicated than just deleting all rows, pass an async `cleanup` option.
For cases where the credentials are not generated with SAS (that is, when testing the Auth service), pass `{noSasCredentials: true}` and the original credentials will be used in "real" mode.

This function assumes the following config values:
 * `cfg.azure.accountId`
 * `cfg.taskcluster.rootUrl`
 * `cfg.taskcluster.credentials`

And assumes that the `exports` argument has a `load` function corresponding to a sticky loader.

Utilities
---------

### Sleep

The `sleep` function returns a promise that resolves after a delay.

**NOTE** tests that depend on timing are notoriously unreliable, and suggest
poorly-isolated tests. Consider writing the tests to use a "fake" clock or to
poll for the expected state.

### Poll

The `poll` function will repeatedly call a function that returns a promise
until the promise is resolved without errors.

```javascript
await poll(
  maybeFunc, // function to be called
  11,        // max times to try it
  100);      // delay (ms) between tries
```
