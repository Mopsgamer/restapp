import type { Config } from "tailwindcss";
const config: Config = {
    content: [
        "./web/templates/**/*.html",
    ],
    theme: {
        extend: {
            transitionDuration: {
                "1100": "1100ms",
                "1200": "1200ms",
                "1300": "1300ms",
                "1400": "1400ms",
            },
            colors: {
                n50: "var(--sl-color-neutral-50)",
                n100: "var(--sl-color-neutral-100)",
                n200: "var(--sl-color-neutral-200)",
                n300: "var(--sl-color-neutral-300)",
                n400: "var(--sl-color-neutral-400)",
                n500: "var(--sl-color-neutral-500)",
                n600: "var(--sl-color-neutral-600)",
                n700: "var(--sl-color-neutral-700)",
                n800: "var(--sl-color-neutral-800)",
                n900: "var(--sl-color-neutral-900)",
            },
        },
    },
};
export default config;
