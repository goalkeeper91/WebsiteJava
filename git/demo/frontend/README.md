# Demo Frontend

Dies ist das Frontend für die Demo-Anwendung. Das Projekt basiert auf React und verwendet Docker zur vereinfachten Entwicklung und Bereitstellung.

## Projektstruktur

- `src/` – Hauptquellcode des Frontends
    - `assets/` – Statische Assets wie Bilder und Icons
    - `components/` – React-Komponenten
        - `admin/` – Admin-spezifische Komponenten
            - `forms/` – Formulare für Admin-Interfaces
            - `tabs/` – Tab-Interfaces
                - `bots/` – Bot-spezifische Tabs
                    - `discord/` – Discord-Bot Komponenten
                - `pages/` – Seiten-Tab-Komponenten
                - `stats/` – Statistik-Tab-Komponenten
        - `community/` – Community-Komponenten
        - `hero/` – Hero-Section Komponenten
        - `live/` – Live-Streaming Komponenten
        - `routes/` – Routing-Komponenten
        - `socials/` – Social-Media-Komponenten
        - `videoComponents/` – Komponenten für Videos

## Setup mit Docker

Das Projekt nutzt Docker, daher ist keine lokale Installation von Node.js oder npm nötig.

### Docker Build & Run

```bash
docker compose up --build -d
```

## React + TypeScript + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type-aware lint rules:

```js
export default tseslint.config([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...

      // Remove tseslint.configs.recommended and replace with this
      ...tseslint.configs.recommendedTypeChecked,
      // Alternatively, use this for stricter rules
      ...tseslint.configs.strictTypeChecked,
      // Optionally, add this for stylistic rules
      ...tseslint.configs.stylisticTypeChecked,

      // Other configs...
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```

You can also install [eslint-plugin-react-x](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'

export default tseslint.config([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs['recommended-typescript'],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```
