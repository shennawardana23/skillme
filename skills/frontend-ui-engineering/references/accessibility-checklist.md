# Accessibility Checklist (WCAG 2.1 AA)

Read this before a final accessibility pass on a non-trivial component or
page — SKILL.md covers the everyday defaults; this is the fuller pass.

## Forms

- Every input has a programmatically associated label (`<label for>`
  matching `id`, or `aria-label`/`aria-labelledby` when no visible label
  exists).
- Required fields are marked both visually and with `aria-required="true"`
  or the native `required` attribute.
- Validation errors are associated with their field via
  `aria-describedby`, not only shown as nearby text.
- Error summaries (for forms with multiple errors) get focus on submit
  failure so screen reader users land on the summary, not wherever focus
  happened to be.
- Fieldsets (`<fieldset>` + `<legend>`) group related inputs like radio
  button sets or address fields.

## Modals and dialogs

- Focus moves into the dialog when it opens (see SKILL.md's focus-move
  example) and returns to the triggering element when it closes.
- Focus is trapped inside the dialog while open — Tab/Shift+Tab cycles
  within the dialog, never escapes to background content.
- `role="dialog"` and `aria-modal="true"` are set; `aria-labelledby` points
  at the dialog's heading.
- `Escape` closes the dialog.
- Background content is inert to screen readers while the dialog is open
  (native `<dialog>` handles this; a custom overlay needs
  `aria-hidden="true"` on siblings or `inert`).

## Live regions and dynamic content

- Content that updates without a page navigation (toast notifications,
  async validation results, a live counter) sits inside an element with
  `aria-live="polite"` (most cases) or `aria-live="assertive"` (urgent,
  rare) so assistive tech announces the change.
- Loading indicators use `aria-busy="true"` on the container, not just a
  visual spinner.
- Don't wrap large regions of frequently-changing content in `aria-live` —
  it causes excessive, unhelpful announcements. Scope it to the specific
  element that changed.

## Keyboard and focus order

- Tab order follows visual/logical reading order — verify by tabbing
  through the page without a mouse.
- No positive `tabindex` values (`tabindex="1"`, `"2"`, ...) — they create a
  separate tab sequence that fights the natural DOM order. Use `0` or `-1`
  only.
- Custom interactive widgets (dropdowns, comboboxes, tab lists) implement
  the arrow-key/Enter/Escape conventions from the WAI-ARIA Authoring
  Practices for that widget pattern, not an ad hoc key scheme.
- Focus is visibly indicated (`:focus-visible` styling) — never
  `outline: none` without a replacement indicator.

## Testing tools

- **axe-core** (via `@axe-core/playwright` or the browser extension) for
  automated WCAG rule checking — catches missing labels, contrast
  failures, and ARIA misuse mechanically.
- **Lighthouse accessibility audit** (Chrome DevTools or CI) for a broader
  pass including heading structure and landmark regions.
- **Manual keyboard pass**: unplug the mouse, Tab through the entire flow,
  confirm every action is reachable and every focus state is visible.
- **Screen reader spot-check**: VoiceOver (macOS/iOS) or NVDA (Windows) on
  the specific component you changed — automated tools don't catch
  everything (e.g., an announcement that's technically valid ARIA but
  confusing when read aloud).
