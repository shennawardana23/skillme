---
name: php-codeigniter-legacy-patterns
description: Guides classic CodeIgniter 2.x/3.x applications (the application/ directory, bcit-ci/CodeIgniter, no PHP namespaces) - a fundamentally different framework generation from CodeIgniter 4. Use when working in a codebase with an application/ directory and CI_VERSION below 4, extending MY_Controller or a *_model class, adding a route, or reviewing code that assumes CI4-style namespaces or DI.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "php"
---

# Classic CodeIgniter (2.x/3.x) Patterns

Classic CodeIgniter (2.x/3.x — recognizable by an `application/` directory,
`bcit-ci/CodeIgniter` in `composer.json`, and no PHP namespaces anywhere in
app code) is not an older version of the same thing this catalog's
`php-codeigniter-patterns`/`-security`/`-tdd`/`-verification` skills teach —
those assume CI4's namespaced classes, `CodeIgniter\Model`, `Config\*`
classes, and the `spark` CLI, none of which exist here. Treat a codebase
with an `application/` directory (not `app/`) as needing this skill
instead, not a CI4 skill applied loosely.

## Everything is global, resolved through a service locator

There is no dependency injection. A controller pulls in what it needs via
`$this->load->`, and that call populates a property on the controller
instance at runtime — nothing is declared as a constructor dependency or
type-checkable ahead of time:

```php
class MY_Controller extends CI_Controller
{
    public function __construct()
    {
        parent::__construct();
        $this->load->model(array('User_model', 'Hotel_model'), '', TRUE);
        $this->load->library('session');
    }
}
```

Every class name here is global — no namespaces — so `User_model` must be
unique across the entire loaded application, and a helper function calling
`get_instance()` (CodeIgniter's way of reaching the current controller
instance from outside it) creates an implicit dependency on whatever
properties that controller's constructor happened to populate:

```php
function requires_admin_privilege() {
    $CI =& get_instance();
    return $CI->privilege_lib->user_is_admin(); // assumes $CI->privilege_lib exists
}
```

A helper like this silently breaks if called from a controller that
doesn't extend the same base class or load the same library — there is no
compiler or type system to catch the mismatch, only a runtime error.

## A shared base controller is usually a god class

It's common for one `MY_Controller` (extended by every controller in the
app, including API controllers) to accumulate constructor logic, loaded
libraries/models, and helper methods for the whole application over time,
because there's no lighter-weight alternative built into the framework —
there is no separate, minimal base for an API-only controller unless the
team builds one deliberately. Before adding "just one more thing" to a
shared base controller, consider whether it actually needs to run for
every controller that extends it, including ones (an API endpoint, a CLI
task) that may not need session/view-rendering setup at all.

## Configuration is a plain PHP array, not typed

```php
// application/config/database.php
$db['default'] = array(
    'hostname' => 'localhost',
    'username' => 'root',
    'password' => '',
    'database' => 'app',
    'dbdriver' => 'mysqli',
);
```

Multiple named connection groups can coexist in the same array (`$db['default']`,
`$db['reporting']`, and so on), selected explicitly per model
(`$this->load->database('reporting')`) — nothing enforces which group a
given model is supposed to use; that convention lives only in
documentation or in the developer's head, not in the type system.

## Routes are manual and exhaustive

`application/config/routes.php` is a flat associative array mapping every
individual URI pattern to a controller/method — there is no automatic
resource-routing convention (`Route::resource`-style) to lean on. A new
CRUD action means an explicit new line in this file; an unmapped route
falls through to the framework's default 404 behavior, and a typo'd
pattern here fails exactly the same silent way a CI4 filter-alias typo
does — the route just doesn't match, with no error pointing at the cause.

## Gotchas

- **No namespaces means every class name is a global collision risk** —
  two libraries or models with the same class name (even in different
  directories) will conflict; check for name collisions explicitly rather
  than assuming directory structure provides isolation the way a
  namespace would.
- **`get_instance()` inside a helper function creates an implicit runtime
  dependency on the calling controller's setup** — a helper that assumes
  `$CI->some_library` exists breaks with no compile-time warning if called
  from a controller that never loaded that library.
- **A shared base controller extended by every controller, API included,
  tends to accumulate unrelated setup logic over time** — evaluate whether
  new base-controller logic genuinely needs to run for every subclass
  before adding it, rather than defaulting to "put it in the base
  controller since everything extends it anyway."
- **Multiple database connection groups in one config file have no
  enforced usage convention** — confirm which group a given model
  actually expects before assuming `$this->db` (the default connection)
  is always the right one.
- **`enable_hooks = TRUE` in `application/config/hooks.php` does not mean
  hooks are actually in use** — check `application/hooks/` for real,
  registered hook points before assuming this is an active extension
  mechanism in a given codebase; it's commonly left enabled with nothing
  registered.
- **A manually-exhaustive `routes.php` has the same silent-typo failure
  mode as a CI4 filter alias typo** — a misspelled pattern doesn't error,
  it just never matches, and the request falls through to a 404 with
  nothing pointing at the routes file as the cause.
- **Vendored, pre-Composer third-party code (a library dropped directly
  into `application/third_party/`) can coexist with Composer-managed
  dependencies in the same app** — these are two independently-maintained
  dependency mechanisms; updating one doesn't touch the other, and it's
  easy to forget the vendored copy exists at all when auditing
  dependencies.

## Real-world grounding

Classic CodeIgniter's lack of namespaces and its `$this->load->`
service-locator style are direct artifacts of the PHP 5.2/5.3-era language
features it targeted — namespaces were only added to PHP in 5.3, and much
of CodeIgniter 2's design predates widespread namespace adoption in the PHP
ecosystem. This is exactly why CodeIgniter 4 (a full rewrite, not an
incremental upgrade) introduced namespaces, PSR-4 autoloading, and a
DI-adjacent `Config`/`Services` model — the two framework generations are
different enough in structure that "upgrade CI3 to CI4" is properly
understood as a migration project, not a version bump, which is the same
reason this catalog treats them as two separate skills rather than one.

## Verification

- [ ] New class names are checked for collisions across the whole
      `application/` tree, since there are no namespaces to isolate them
- [ ] Helper functions using `get_instance()` document which
      controller-populated properties they depend on
- [ ] New logic added to a shared base controller genuinely needs to run
      for every subclass, including any API-only controllers that extend
      it
- [ ] The correct database connection group is used for a given model,
      confirmed rather than defaulted to `$this->db`
- [ ] `application/hooks/` actually has registered hooks before treating
      hooks as an active extension mechanism in this codebase
- [ ] A new route added to `routes.php` has been tested by an actual
      request, not just read for correctness
