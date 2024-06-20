/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  darkMode: "class",
  theme: {
    extend: {
      transitionProperty: {
        height: "height",
        width: "width",
        padding: "padding",
      },
      colors: {
        primary: "#3EB2B2",
      },
      backgroundColor: {
        primary: "#3EB2B2",
        menu: "#F3F3F3",
      },
      fontFamily: {
        inter: ["Inter", "sans-serif"],
      },
    },
  },
  plugins: [],
};
