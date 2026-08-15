---
name: vue-nuxt-frontend-patterns
description: Guides Vue 3 Composition API and Nuxt 3 development — composables, reactivity, SSR data fetching, and component patterns. Use when building Vue or Nuxt components, writing composables, fetching data in a Nuxt page or component, managing shared state with Pinia, or debugging a value that "isn't reactive" after being passed around.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Vue 3 / Nuxt 3 Frontend Patterns

Composition API over Options API for new code. Reactivity is opt-in per
value (`ref`/`reactive`), not automatic on plain destructuring — most
"reactivity broke" bugs come from losing that connection, not from Vue
itself misbehaving.

## The reactivity-loss gotcha (highest-value thing to know)

Destructuring a `reactive()` object breaks reactivity on the destructured
values — they become plain, disconnected values at the moment of
destructuring:

```typescript
const state = reactive({ count: 0, name: "task" });

// BROKEN: count is now a plain number, frozen at 0 forever
const { count } = state;

// CORRECT: toRefs preserves the reactive connection
const { count } = toRefs(state);
```

`ref()` values avoid this entirely because the ref object itself is passed
around and only unwrapped at the point of use — prefer `ref()` for
single values, `reactive()` only for grouped object state you won't
destructure before the template.

## Composables

A composable is a function starting with `use` that encapsulates
reusable reactive logic — Vue's equivalent of a custom React hook, but
without the rules-of-hooks ordering constraint, since Vue's reactivity
doesn't depend on call order.

```typescript
// composables/useDebounce.ts
export function useDebounce<T>(source: Ref<T>, delayMs: number) {
  const debounced = ref(source.value) as Ref<T>;
  let timer: ReturnType<typeof setTimeout>;
  watch(source, (value) => {
    clearTimeout(timer);
    timer = setTimeout(() => { debounced.value = value; }, delayMs);
  });
  return debounced;
}
```

## Nuxt 3 data fetching

Use `useFetch`/`useAsyncData` inside `<script setup>` for SSR-safe data
fetching — they dedupe requests between server and client render and
integrate with Nuxt's payload serialization. Use plain `$fetch` only inside
event handlers or outside the SSR render path (a click handler, a
composable called imperatively, not during initial setup).

```vue
<script setup lang="ts">
const route = useRoute();
const { data: reservation, error, pending } = await useFetch(
  () => `/api/reservations/${route.params.id}`
);
</script>
```

## Gotchas

- Calling `$fetch` directly inside `<script setup>`'s top-level body (rather
  than `useFetch`/`useAsyncData`) runs the request twice — once on the
  server during SSR, once again on the client during hydration — because
  there's no dedup/payload mechanism attached to a bare `$fetch` call there.
- Mutating a prop directly (`props.items.push(x)`) works in dev but creates
  a component whose behavior depends on the parent's internal object
  identity — treat props as read-only; emit an event or use `v-model` for
  two-way binding instead.
- A `v-for` without a stable `:key` (using the array index when items can
  reorder or be removed) causes Vue to reuse the wrong DOM node's local
  state across re-renders — a component that "remembers" input from a
  different row after a delete is almost always this.
- Server-rendered output that differs from what the client re-renders
  during hydration (e.g., using `Date.now()` or `Math.random()` directly in
  a template) produces a hydration mismatch — this class of bug is common
  to every SSR framework (Vue, Nuxt, Next, SvelteKit), not unique to Nuxt,
  and the fix is always the same: make the value deterministic between
  server and client, or defer it to a client-only lifecycle hook.

## State sharing (Pinia)

```typescript
export const useReservationStore = defineStore("reservations", () => {
  const items = ref<Reservation[]>([]);
  const confirmed = computed(() => items.value.filter(r => r.status === "confirmed"));
  async function load(hotelId: string) {
    items.value = await $fetch(`/api/hotels/${hotelId}/reservations`);
  }
  return { items, confirmed, load };
});
```

## Verification

- [ ] No `reactive()` object destructured without `toRefs`
- [ ] Data fetching in `<script setup>` uses `useFetch`/`useAsyncData`, not a bare top-level `$fetch`
- [ ] Props are never mutated directly
- [ ] Every `v-for` has a stable, unique `:key`
- [ ] No `Date.now()`/`Math.random()`/browser-only API used directly in a template that renders during SSR
