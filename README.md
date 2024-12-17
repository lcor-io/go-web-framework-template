Template for a Site Generation Framework in Go, with as little as JavaScript as possible. If one doesn't need Tailwind for css, you can even remove JS from the build steps entirely.

 - Web Framework - [Fiber](https://gofiber.io/)
 - HTML generation - [Templ](https://templ.guide/)
 - CSS - [TailwindCSS](https://tailwindcss.com/)
 - Interactivity - [htmx](https://htmx.org/)
 
As of now, it supports Static Page Generation and Dynamic Rendering with two different renderers.

## First install

    go mod tidy
    pnpm install  # Or whatever package manager you are using

## Launch
Hot reloading should launch automatically in the browser.

    make live
    
To cleanup build files after dev

    make clean

## TODOS

 - [x] ~~More route examples, with both SSG and Dynamic rendering~~
 - [ ] Example of a Templ component using htmx
 - [ ] Initial properties at build time
 - [ ] Dynamic routing system using FS ?
 - [ ] Persistent cache files between restarts. Dig into Fiber cache middleware
 - [ ] Layout composition on children
 - [ ] i18n
 - [ ] Add template to be "PWA ready"
