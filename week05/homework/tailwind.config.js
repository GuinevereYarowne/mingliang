/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}", // 匹配src下所有ts/tsx/js/jsx文件
  ],
  theme: {
    extend: {
      colors: {
        darkBlue: "#0F172A",
        cardBlue: "#1E293B",
      },
    },
  },
  plugins: [],
};