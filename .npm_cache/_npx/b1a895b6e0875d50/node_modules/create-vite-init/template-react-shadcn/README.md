# React + TypeScript + Vite + Shadcn UI 

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

## Technical Stack & Constraints
- **Core Technologies**: 
  ```text
  React 18, TypeScript, Vite 3, Tailwind CSS 3, shadcn/ui, React Router 6
  ```
- **Project Structure**:
  - All code under `src/`
  - Critical Files:
    - `main.tsx` (read-only entry, only contains `render(<App />)`, DO NOT generate)
    - `App.tsx` (modifiable, define global providers here, export as default)
    - `router.tsx` (centralized routing configuration)
  - Directories:
    - `components/` (reusable components, with `ui/` for shadcn components via `@/components/ui` alias)
    - `pages/` (page components mapped to routes)
    - `layouts/` (shared layout components)
    - `lib/` (utilities/hooks)
    - `assets/` (static files)

## Routing Requirements
1. **Multi-Level Routing**:
   - Use `createBrowserRouter` + `RouterProvider` (React Router 6 syntax)
   - Example hierarchy:
     ```tsx
     // router.tsx
     const router = createBrowserRouter([
       {
         path: "/",
         element: <MainLayout />, // Contains nav/footer
         children: [
           { path: "/", element: <HomePage /> },
           { path: "/products", element: <ProductLayout />, 
             children: [
               { path: ":id", element: <ProductDetail /> }
             ]
           }
         ]
       }
     ])
     ```
   - Use `Outlet` for nested layouts
   - **Strict Avoidance**: No `window.location` or `<a>` tags - use `Link`/`Navigate`

## UI/UX Specifications
### Styling Guidelines
- **Tailwind Best Practices**:
  - Use `@layer` for custom utilities
  - Responsive breakpoints:
    ```css
    /* Mobile-first approach */
    md: (768px), lg: (1024px), xl: (1280px)
    ```
  - Theme colors from `tailwind.config.js`
- **Shadcn UI**:
  - Available components:
    ```text
    accordion, alert, alert-dialog, aspect-ratio, avatar, badge, breadcrumb, button, calendar, card, carousel, chart, checkbox, collapsible, command, context-menu, dialog, drawer, dropdown-menu, form, hover-card, input, input-otp, label, menubar, navigation-menu, pagination, popover, progress, radio-group, resizable, scroll-area, select, separator, sheet, sidebar, skeleton, slider, sonner, switch, table, tabs, textarea, toggle, toggle-group, tooltip
    ```
- **Animation**:
  - Transition examples:
    ```tsx
    <Card className="transition-all hover:scale-105 duration-300">
    ```

## Quality Assurance
1. **Code Standards**:
   - TypeScript strict mode
   - ESLint/Prettier preconfigured
2. **Performance**:
   - Route-based code splitting (`lazy()` + `Suspense`)
3. **Error Handling**:
   - React Router `errorElement` boundaries
   - Global error fallback UI

## Special Constraints
- **Shadcn Overrides**:
  ```text
  Customize via Tailwind - do NOT edit shadcn/core CSS directly
  ```
- **Asset Handling**:
  ```text
  Images: Use `src/assets/` with Vite import syntax
  Icons: Use `@radix-ui/react-icons`
  ```
