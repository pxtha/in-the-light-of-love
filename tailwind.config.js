/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.html"],
  theme: {
    extend: {
      colors: {
        'velvet-red': '#7A1C1C',
        'bright-red': '#C41E3A',
        'candlelight-yellow': '#FFD67A',
        'deep-black': '#0D0D0D',
        'ivory-white': '#F5F5F5',
      },
      fontFamily: {
        'serif': ['"Playfair Display"', 'serif'],
        'sans': ['"Inter"', 'sans-serif'],
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        lightoflove: {
          "primary": "#7A1C1C",
          "secondary": "#C41E3A",
          "accent": "#FFD67A",
          "neutral": "#0D0D0D",
          "base-100": "#F5F5F5",
          "info": "#3ABFF8",
          "success": "#36D399",
          "warning": "#FBBD23",
          "error": "#F87272",
        },
      },
    ],
  },
}
