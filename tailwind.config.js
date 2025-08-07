/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./frontend/**/*.templ",
  ],
  safelist: [
    {
        pattern: /.*/,
    }
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}