/**
 * @typedef {import('./utils.js').GameSettings} GameSettings
 */

import {
    assertNotNull,
    isAlive,
    parseGameSettings,
    parseNumber,
    validatePositiveInteger
} from "./utils.js";

const baseURL = "/api/v1/game"

export class GameOfLifeViewer extends HTMLElement {
    static observedAttributes = ["rows", "init-state", "seed"];

    /** @type {GameSettings|null} */
    #settings = null

    /** @type {EventSource|null} */
    #events = null

    /** @type {HTMLCanvasElement} */
    #canvasElement;

    /** @type {ShadowRoot} */
    #shadow;

    /** @type {number|null} */
    #lastRAFid = null

    constructor() {
        super();
        this.#canvasElement = document.createElement("canvas")
        this.#canvasElement.style.backgroundColor="black"
        this.#shadow = this.attachShadow({mode: "closed" })
    }

    /**
     * @returns {number|null}
     */
    get rows() {
        const raw = this.getAttribute("rows")
        if (raw == null) {
            return null;
        }

        const rows = parseNumber(raw);
        validatePositiveInteger(rows);
        return rows
    }

    /**
     * @returns {string|null}
     */
    get initState() {
        return this.getAttribute("init-state")
    }

    /**
     * @returns {number|null}
     */
    get seed() {
        const raw = this.getAttribute("seed")
        if (raw == null) {
            return null;
        }

        const seed = parseNumber(raw);
        validatePositiveInteger(seed);
        return seed
    }

    get #searchParams() {
        const params = new URLSearchParams()
        if (this.seed != null) {
            params.set("seed", String(this.seed))
        }
        if (this.initState != null) {
            params.set("init_state", this.initState)
        }
        if (this.rows != null) {
            params.set("rows", String(this.rows))
        }
        return params
    }

    // noinspection JSUnusedGlobalSymbols
    connectedCallback() {
        const eventsURL = new URL(baseURL, location.origin);
        eventsURL.search = this.#searchParams.toString();
        this.#events = new EventSource(eventsURL);

        this.#events.addEventListener("error", (event)=>{
            console.error("Failed to connect to events source", event)
        })
        this.#events.addEventListener("open", ()=>{
            this.#shadow.appendChild(this.#canvasElement)
        })
        this.#events.addEventListener("settings", this.#handleSettings.bind(this))
        this.#events.addEventListener("message", this.#handleMessage.bind(this))
    }

    /**
     * @param {MessageEvent} event
     */
    #handleSettings(event) {
        this.#settings = parseGameSettings(event.data)
        this.#canvasElement.width = this.#settings.boardWidth
        this.#canvasElement.height = this.#settings.boardHeight
        this.#canvasElement.dataset.rows = this.#settings.rows.toString()
        this.#canvasElement.dataset.columns = this.#settings.columns.toString()
    }

    /**
     * @param {MessageEvent} event
     */
    #handleMessage(event) {
        this.#cancelAnimationFrameRequest()

        assertNotNull(this.#settings)
        const {rows, columns, cellWidth, cellHeight} = this.#settings
        const binData = /** @type {ArrayLike<string>} */atob(event.data)
        const u8arr = Uint8Array.from(binData, b64byte => b64byte.codePointAt(0))
        const padding = u8arr[0]
        const data = u8arr.slice(1)

        this.#lastRAFid = requestAnimationFrame(() => {
            for (let row = 0; row < rows; row++) {
                for (let column = 0; column < columns; column++) {
                    const ctx = this.#canvasElement.getContext("2d")
                    assertNotNull(ctx)

                    ctx.fillStyle = "white";
                    if (isAlive(this.#settings, data, {row, column})) {
                        ctx.fillRect(column * cellWidth, row * cellHeight, cellWidth, cellHeight)
                    } else {
                        ctx.clearRect(column * cellWidth, row * cellHeight, cellWidth, cellHeight)
                    }
                }
            }
        })
    }

    #cancelAnimationFrameRequest() {
        if (this.#lastRAFid != null) {
            cancelAnimationFrame(this.#lastRAFid)
        }
    }

    // noinspection JSUnusedGlobalSymbols
    disconnectedCallback() {
        this.#events.close()
        this.#cancelAnimationFrameRequest()
    }

    // noinspection JSUnusedGlobalSymbols
    attributeChangedCallback() {
        this.disconnectedCallback()
        this.connectedCallback()
    }

}

window.customElements.define("game-of-life-viewer", GameOfLifeViewer);