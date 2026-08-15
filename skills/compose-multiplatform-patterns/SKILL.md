---
name: compose-multiplatform-patterns
description: Patterns for building shared UI across Android, iOS, Desktop, and Web with Compose Multiplatform and Jetpack Compose — state hoisting with ViewModel/StateFlow, type-safe navigation, slot-based composables, recomposition performance, and expect/actual platform code. Use when writing or reviewing Compose UI, wiring a ViewModel to a screen, debugging unnecessary recomposition, or adding platform-specific behavior in a KMP project.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Compose Multiplatform Patterns

Shared UI across Android, iOS, Desktop, and Web via Compose Multiplatform,
plus Jetpack Compose specifics on Android. Core mental model: state flows
down as an immutable snapshot, events flow up as callbacks — the same
unidirectional-data-flow discipline as SwiftUI or React, expressed with
Kotlin's `StateFlow` and Compose's recomposition system instead of a virtual
DOM diff.

## When to use

Building Compose UI (Jetpack Compose or Compose Multiplatform); wiring
ViewModels to screen state; implementing navigation in a KMP or Android
project; designing reusable composables; diagnosing excess recomposition or
janky lists.

## State management

**Single state data class per screen**, exposed as `StateFlow`, collected
with lifecycle awareness:

```kotlin
data class ItemListState(
    val items: List<Item> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null,
    val searchQuery: String = ""
)

class ItemListViewModel(private val getItems: GetItemsUseCase) : ViewModel() {
    private val _state = MutableStateFlow(ItemListState())
    val state: StateFlow<ItemListState> = _state.asStateFlow()

    fun onSearch(query: String) {
        _state.update { it.copy(searchQuery = query) }
        viewModelScope.launch {
            _state.update { it.copy(isLoading = true) }
            getItems(query).fold(
                onSuccess = { items -> _state.update { it.copy(items = items, isLoading = false) } },
                onFailure = { e -> _state.update { it.copy(error = e.message, isLoading = false) } }
            )
        }
    }
}

@Composable
fun ItemListScreen(viewModel: ItemListViewModel = koinViewModel()) {
    val state by viewModel.state.collectAsStateWithLifecycle()
    ItemListContent(state = state, onSearch = viewModel::onSearch)
}

@Composable
private fun ItemListContent(state: ItemListState, onSearch: (String) -> Unit) {
    // Stateless — trivially previewable and testable in isolation
}
```

`collectAsStateWithLifecycle` (not plain `collectAsState`) stops collecting
when the composable leaves the started lifecycle state, avoiding wasted work
and leaked collection in backgrounded screens.

**Event sink for complex screens**: replace a growing pile of callback
lambdas with one sealed-interface event and a single `onEvent` lambda —
easier to log, test, and extend without changing every call site:

```kotlin
sealed interface ItemListEvent {
    data class Search(val query: String) : ItemListEvent
    data class Delete(val itemId: String) : ItemListEvent
    data object Refresh : ItemListEvent
}

fun onEvent(event: ItemListEvent) = when (event) {
    is ItemListEvent.Search -> onSearch(event.query)
    is ItemListEvent.Delete -> deleteItem(event.itemId)
    is ItemListEvent.Refresh -> loadItems(_state.value.searchQuery)
}
```

## Navigation

Compose Navigation 2.8+ supports type-safe routes as `@Serializable`
objects/classes instead of string paths — the compiler catches a missing or
mistyped argument instead of a runtime crash on navigation:

```kotlin
@Serializable data object HomeRoute
@Serializable data class DetailRoute(val id: String)

@Composable
fun AppNavHost(navController: NavHostController = rememberNavController()) {
    NavHost(navController, startDestination = HomeRoute) {
        composable<HomeRoute> {
            HomeScreen(onNavigateToDetail = { id -> navController.navigate(DetailRoute(id)) })
        }
        composable<DetailRoute> { backStackEntry ->
            DetailScreen(id = backStackEntry.toRoute<DetailRoute>().id)
        }
    }
}
```

Use `dialog<Route>()` for dialogs/bottom sheets so open/close state lives in
the back stack, not in an imperative `showDialog` boolean scattered across
composables.

## Composable design

**Slot-based APIs** over one big prop bag — callers compose exactly the
content they need instead of fighting a fixed configuration surface:

```kotlin
@Composable
fun AppCard(
    modifier: Modifier = Modifier,
    header: @Composable () -> Unit = {},
    content: @Composable ColumnScope.() -> Unit,
    actions: @Composable RowScope.() -> Unit = {}
) {
    Card(modifier = modifier) {
        Column {
            header()
            Column(content = content)
            Row(horizontalArrangement = Arrangement.End, content = actions)
        }
    }
}
```

**Modifier order changes rendered output** — modifiers apply as a chain, so
`padding().background()` and `background().padding()` produce visibly
different results:

```kotlin
Modifier
    .padding(16.dp)                    // 1. layout: reserve space first
    .clip(RoundedCornerShape(8.dp))    // 2. shape the drawing area
    .background(Color.White)          // 3. draw within the clipped shape
    .clickable { }                    // 4. attach interaction last
```

## Platform-specific code (expect/actual)

```kotlin
// commonMain
@Composable
expect fun PlatformStatusBar(darkIcons: Boolean)

// androidMain
@Composable
actual fun PlatformStatusBar(darkIcons: Boolean) {
    val systemUiController = rememberSystemUiController()
    SideEffect { systemUiController.setStatusBarColor(Color.Transparent, darkIcons) }
}
```

Keep `expect` declarations at the smallest surface that actually differs per
platform — a whole differing screen belongs behind a shared interface
injected per-platform, not an `expect fun` per widget.

## Performance

**Stability drives skippability.** Compose skips recomposing a function
whose inputs are unchanged only if it can prove those inputs are *stable*.
Mark data classes `@Immutable` (all properties never change after
construction) or `@Stable` (properties may change but notify Compose via
`equals`/snapshot state) so the compiler can actually make that guarantee
instead of conservatively recomposing every time:

```kotlin
@Immutable
data class ItemUiModel(val id: String, val title: String, val progress: Float)
```

**Stable keys in lazy lists** let Compose track which physical row maps to
which data item across reordering, insertion, and deletion, so item
identity — and any `remember`ed per-item state — survives instead of being
reset:

```kotlin
LazyColumn {
    items(items = items, key = { it.id }) { item -> ItemRow(item = item) }
}
```

**Defer expensive reads with `derivedStateOf`** so a value that's expensive
to compute — or that only needs to change occasionally — doesn't trigger
recomposition on every upstream state tick:

```kotlin
val listState = rememberLazyListState()
val showScrollToTop by remember {
    derivedStateOf { listState.firstVisibleItemIndex > 5 }
}
```

**Don't allocate a fresh lambda per item inside a loop without a stable
key** — Compose can't tell the new lambda apart from the old one across
recompositions, so any state tied to that row's identity (focus, animation)
resets. Fixing this is two separate concerns: hoist the filter/derived list
out with `remember(items)` so it isn't recomputed every recomposition, *and*
wrap each row in `key(item.id) { ... }` so Compose tracks row identity
independent of the lambda captured inside it.

## Theming

```kotlin
@Composable
fun AppTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    dynamicColor: Boolean = true,
    content: @Composable () -> Unit
) {
    val colorScheme = when {
        dynamicColor && Build.VERSION.SDK_INT >= Build.VERSION_CODES.S ->
            if (darkTheme) dynamicDarkColorScheme(LocalContext.current)
            else dynamicLightColorScheme(LocalContext.current)
        darkTheme -> darkColorScheme()
        else -> lightColorScheme()
    }
    MaterialTheme(colorScheme = colorScheme, content = content)
}
```

## Gotchas

- `mutableStateOf` inside a ViewModel survives configuration changes but
  isn't lifecycle-aware the way `StateFlow` + `collectAsStateWithLifecycle`
  is — prefer `StateFlow` for anything a screen collects.
- Passing `NavController` deep into child composables couples them to
  navigation; pass lambda callbacks instead so a composable can be reused or
  previewed without a real `NavHost`.
- Heavy computation inside a `@Composable` function body reruns on every
  recomposition unless wrapped in `remember{}` — move it to the ViewModel or
  memoize it explicitly.
- `LaunchedEffect(Unit)` is not a ViewModel-init substitute: on some
  navigation/configuration-change paths it re-runs, silently re-triggering
  "one-time" setup.
- Creating a new object instance (a new lambda, a new data class instance)
  as a composable parameter on every call defeats stability checks even if
  the underlying data hasn't changed — construct it once and hoist it.

## Real-world grounding

The `@Immutable`/`@Stable` contract and its effect on skippable
recomposition is Jetpack Compose's own documented performance model — this
is not a project convention, it's how the compiler and runtime decide
whether to skip a composable. Compose Multiplatform (JetBrains' extension of
Jetpack Compose to iOS, Desktop, and Web) inherits this same recomposition
model, so a stability bug reproduces identically across platforms even
though the surrounding platform code differs.
