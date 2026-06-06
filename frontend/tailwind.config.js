export default {
    content: ["./index.html", "./src/**/*.{ts,tsx}"],
    theme: {
        extend: {
            colors: {
                surface: "#111827",
                "surface-2": "#1F2937",
                accent: "#3B82F6",
            },
            fontFamily: {
                mono: ["JetBrains Mono", "monospace"],
                sans: ["Inter", "sans-serif"],
            },
        },
    },
    plugins: [],
};
