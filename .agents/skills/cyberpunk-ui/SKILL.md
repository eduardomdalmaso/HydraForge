---
name: cyberpunk-ui
description: Design system guidelines, tokens, UI components, animations, alerts, and modular CSS rules (<100 lines per file) for Cyberpunk High-Tech HUD interfaces in HydraForge.
---

# ⚡ Cyberpunk High-Tech HUD Design System & UI Skill

This skill defines the complete visual aesthetic, design tokens, UI components, and architectural rules for HydraForge and HydraStream interfaces.

---

## 1. Aesthetic Philosophy (Night City 2077 HUD)

HydraForge is a high-performance computer vision cockpit engineered for the NVIDIA RTX 5090 and Ultralytics YOLO26. The interface must look and feel like an authentic high-density Cyberpunk HUD:
- **Dark Optical Vacuum:** Deep void background (`#07080c`) with subtle cyan gridlines.
- **Neon Laser Accents:** Intense glows (`#00f0ff`, `#fcee0a`, `#ff0055`, `#00ff66`).
- **Tactical Geometry:** Beveled corners (`clip-path: polygon(...)`), HUD corner brackets, targeting reticles, and scanlines.
- **Micro-Animations:** Hover lift, neon pulse, glitch titles, smooth Bézier curves.

---

## 2. Color Palette & Design Tokens

Always use CSS custom properties defined in `web/src/css/variables.css`:

| Token | Hex / Value | Semantic Role |
| :--- | :--- | :--- |
| `--cb-bg-void` | `#07080c` | Deep void backdrop / root body |
| `--cb-bg-surface` | `#0e1017` | Standard card surface |
| `--cb-bg-elevated` | `#161924` | Input boxes, selects, inner panels |
| `--cb-cyan` | `#00f0ff` | Primary neon glow, links, detection reticles, zero-copy metrics |
| `--cb-yellow` | `#fcee0a` | Active items, highlights, card titles, warnings |
| `--cb-magenta` | `#ff0055` | Danger, purge actions, critical alerts, errors |
| `--cb-green` | `#00ff66` | Online status, completed jobs, TensorRT speedups |
| `--cb-border-dim` | `#232736` | Card borders, subtle separators |

---

## 3. Typography Stacks

- **Display & Titles:** `'Oxanium', sans-serif` (headers, primary action buttons).
- **Sub-headings:** `'Tomorrow', sans-serif` (section subtitles).
- **Telemetry & Numbers:** `'JetBrains Mono', monospace` (loss, FPS, bounding box coordinates).
- **Body & Captions:** `'Barlow', sans-serif` (readable descriptions).

---

## 4. Universal Component Patterns

### 4.1. Cyber Alerts & Confirmation Modals
> [!IMPORTANT]
> **NEVER use browser native `window.alert()` or `window.confirm()`.** Generic browser dialogs break the cyberpunk immersion.

Always use `<CyberAlertModal />` or styles from `web/src/css/components/alerts.css`:
```jsx
<CyberAlertModal
  isOpen={isOpen}
  title="SYSTEM CONFIRMATION"
  type="danger" // 'danger' | 'yellow' | 'cyan'
  message="Are you sure you want to execute this operation?"
  confirmText="CONFIRM"
  cancelText="ABORT"
  onConfirm={handleConfirm}
  onClose={() => setIsOpen(false)}
/>
```

Features:
- Glassmorphic backdrop (`backdrop-filter: blur(8px)`).
- Neon glow (`box-shadow: 0 0 35px rgba(...)`).
- HUD targeting corners (`.hud-corner-tl`, `.hud-corner-br`).
- Beveled action buttons (`.cyber-action-btn`).

### 4.2. Detection Menus & Vision Task HUD Reticles
Vision tasks (Detection, Segmentation, Pose, OBB) and detection viewports must use HUD targeting markers:
- `.scanline-overlay`: Subtle scanline scan texture over viewports and canvases.
- `.hud-corner-tl` & `.hud-corner-br`: Yellow/cyan corner brackets.
- `.detection-box` & `.detection-tag`: Neon Cyan hoverable box that transitions to Cyber Yellow on hover.
- `.cyber-pill`: High-tech toggle buttons with hover glow and active cyan luminescence.

### 4.3. Beveled Action Buttons (`.cyber-action-btn`)
```css
.cyber-action-btn {
  background: var(--cb-yellow);
  color: #07080c;
  font-family: var(--font-oxanium);
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 8px), calc(100% - 8px) 100%, 0 100%);
  filter: drop-shadow(-2px -2px 0px var(--cb-cyan)) drop-shadow(2px 2px 0px var(--cb-magenta));
}
```

---

## 5. Strict 100-Line Modularity Rule

> [!WARNING]
> **No `.css`, `.js`, or `.jsx` file in `web/` may exceed 100 lines.**

If a component or style sheet grows beyond 100 lines:
1. Split into focused sub-components (e.g., `CyberAlertModal.jsx`, `DatasetListCard.jsx`).
2. Move helper logic to `web/src/utils/`.
3. Modularize CSS into view-specific or component-specific files in `web/src/css/components/` and `web/src/css/views/`.

---

## 6. Self-Healing & Validation Checklist

Before finalizing any UI changes, execute:
```bash
# 1. Audit line counts (< 100 lines per file)
find web/src -type f \( -name "*.jsx" -o -name "*.js" -o -name "*.css" \) -exec wc -l {} + | grep -v total | awk '$1 > 100'

# 2. Check build
npm --prefix web run build
```
