/**
 * @typedef {import('./app.js').GameOfLifeViewer} GameOfLifeViewer
 */

document.addEventListener("DOMContentLoaded", function () {
    const searchParams = new URLSearchParams(location.search)
    /** @type {GameOfLifeViewer} */
    for (const viewer of document.getElementsByTagName("game-of-life-viewer")) {
        const rows = searchParams.get("rows")
        if (rows != null) {
            viewer.setAttribute("rows", String(searchParams.get("rows")))
        }

        const initState = searchParams.get("init_state")
        if (initState != null) {
            viewer.setAttribute("init-state", String(searchParams.get("init_state")))
        }

        const seed = searchParams.get("seed")
        if (seed != null) {
            viewer.setAttribute("init-state", String(searchParams.get("seed")))
        }
    }
})
