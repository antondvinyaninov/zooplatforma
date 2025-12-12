/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        'vk-blue': '#2787f5',
        'vk-blue-hover': '#1976e3',
        'vk-gray': '#99a2ad',
        'vk-border': '#e7e8ec',
        'vk-bg': '#f0f2f5',
      },
    },
  },
  plugins: [],
}
