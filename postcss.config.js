/** @type {import('postcss-load-config').Config} */
const config = {
  plugins: [
    require("tailwindcss"),
    require("autoprefixer"),
    require("cssnano")({ preset: "default" }),
  ],
};

module.exports = config;
