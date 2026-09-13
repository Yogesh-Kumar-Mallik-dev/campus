# Rule: Mobile-First Adaptation & Viewport Resilience Architecture

## 1. Core Responsive Philosophy: Fluid Layouts & Viewport Resilience

Do NOT approach responsiveness as "add more Tailwind/CSS breakpoints for specific devices."

The target is NOT: "Support every device model."
The target IS: **"Every component should remain usable and visually correct for any reasonable viewport or container size."**

### Architecture Principles:

1. **Fluid Layouts:** Use `w-full`, `max-w-[min(..., 100%)]`, dynamic clamps (`clamp(1rem, 2.5vw, 2rem)`).
2. **Container-Aware Components:** Auto-fit grids (`grid-cols-[repeat(auto-fit,minmax(min(100%,280px),1fr))]`).
3. **Small Number of Structural Breakpoints:**
   - `sm: 640px`
   - `md: 768px` (Primary desktop vs mobile structural shift)
   - `lg: 1024px`
   - `xl: 1280px`
4. **Content-Driven Sizing:** Avoid hardcoded pixel dimensions on text containers and view wrappers.
5. **Defensive Overflow Protection:** Always include `min-w-0`, `flex-wrap gap-2`, and `truncate` where text might expand.
6. **Zero Root Horizontal Overflow:** Page-level horizontal scrolling (`scrollWidth > innerWidth`) is strictly forbidden on root viewports.
7. **Isolated Horizontal Scrolling:** If wide tables or matrices require horizontal space, isolate scrolling strictly inside bounded component wrappers (`overflow-x-auto w-full min-w-0`).
8. **Viewport-Safe Modals & Dialogs:** Constrain all dialogs and bottom sheets to `max-h-[min(90dvh,800px)] overflow-y-auto` with internal scrolling so headers and action buttons remain visible on short screens (e.g. 1280×600 laptop or 844×390 mobile landscape).
9. **Safe Area Insets:** Support notch, home indicator, and foldable displays via `viewport-fit=cover` and safe area insets.

---

## 2. Mobile-First Adaptation: Do NOT Force Desktop UI onto Mobile

Do NOT interpret "responsive" as:

> _"Take the desktop layout and squeeze everything until it fits on mobile."_

That is strictly forbidden.

When a desktop interaction pattern becomes unsuitable for a small screen, use the **appropriate mobile-specific interaction pattern instead**.

### Structural Interaction Mappings

| Desktop UI Interaction | Mobile UI Interaction Pattern (< 768px / md) |
| :--- | :--- |
| **Persistent Top / Side Navigation** | **Hamburger `[☰]` + Slide-Over Drawer / Sheet** containing telemetry, workspaces, and navigation sub-links. |
| **Multi-Item Sub-Header Tabs** | **Breadcrumb + Mobile Section Dropdown (`Select`)** for 1-tap switching without horizontal scrolling. |
| **Wide Tabular Grid (`<Table>`)** | **Dedicated Mobile Entity Card List** (with prominent title, metadata chips, live calculated badges, and full-width touch actions) + optional table view toggle. |
| **Horizontal Toolbar** | **Prominent Search Bar** full-width + Compact Action row + primary `+ Create` button. |
| **Multi-Column Form Grid** | **Single-Column Stacked Form** with generous vertical touch spacing (min 44px touch targets). |
| **Split-Panel / Dual-Axis Inspector** | **Mobile Tabbed Inspector** allowing full-height focused view on both axes. |
| **Horizontal Action Button Trays** | **Stacked Primary Action** (full-width `[Save Changes]`) above secondary actions (`[Delete]`, `[Cancel]`). |

---

## 3. Strict Prohibitions

Do NOT:

- ❌ Shrink everything until it fits
- ❌ Reduce text font sizes until they become illegible (< 11px)
- ❌ Make buttons tiny or smaller than 36px–44px touch targets
- ❌ Squeeze navigation items into a single tiny, overflowing row
- ❌ Force desktop multi-column sidebars onto phones
- ❌ Force desktop wide tables onto narrow mobile screens without a dedicated card alternative
- ❌ Cram toolbars into an unreadable single line
- ❌ Hide critical functionality or primary actions simply to preserve the desktop layout

---

## 4. Architectural Preference Order

When a component does not fit on a smaller screen, follow this strict preference order:

1. **Fluidly resize it** if it remains usable
2. **Reflow / wrap it** if that remains usable
3. **Change the layout structure** (e.g., 2-column to 1-column stack)
4. **Replace desktop interaction with a mobile-specific interaction** (e.g., table $\to$ touch cards; split panels $\to$ mobile tabs)
5. **Move secondary actions into an overflow menu** (`[⋮]` or sheet)
6. **Collapse navigation into a drawer / sheet**
7. **Stack content vertically**
8. **Use isolated horizontal scrolling ONLY when the content genuinely requires it** (e.g., code diffs, wide matrices)

---

## 5. Standardized Pagination & Layout Jump Prevention

1. **Standard 10-Item Page Size:** All collection GET endpoints and data registries must be paginated into standard 10-item pages.
2. **Top Pagination Header Bar:** Position pagination controls (`Showing range`, `[Previous]`, `Page X / Y`, `[Next]`) **ABOVE** the data container. Placing pagination above data ensures users navigate without the controls jumping up and down dynamically based on varying record heights.
3. **Zero Layout Shifts:** Heights and pagination boundaries must be deterministic.
