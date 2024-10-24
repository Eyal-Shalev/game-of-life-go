/**
 * @typedef GameSettings
 * @prop {number} rows
 * @prop {number} columns
 * @prop {number} interval
 * @prop {number} boardHeight
 * @prop {number} boardWidth
 * @prop {number} cellHeight
 * @prop {number} cellWidth
 */

/** @type {GameSettings|null} */
let gameSettings = null

document.addEventListener("DOMContentLoaded", function () {
    const gameUrl = new URL("/api/v1/game", location.origin)
    gameUrl.search = location.search

    const es = new EventSource(gameUrl)

    /** @type HTMLCanvasElement|null */
    const theCanvas = document.querySelector("#the-canvas")
    assertNotNull(theCanvas, "#the-canvas element is missing")

    es.addEventListener("settings", handleSettings(theCanvas))
    es.addEventListener("message", handleMessage(theCanvas))
})

/**
 * @param {HTMLCanvasElement} theCanvas
 * @return {(event: MessageEvent) => void}
 */
function handleSettings(theCanvas) {
    return function (event) {
        gameSettings = parseGameSettings(event.data)
        theCanvas.width = gameSettings.boardWidth
        theCanvas.height = gameSettings.boardHeight
        theCanvas.dataset.rows = gameSettings.rows.toString()
        theCanvas.dataset.columns = gameSettings.columns.toString()
    }
}

/**
 * @type {number|null}
 */
let lastRAFid = null

/**
 * @param {HTMLCanvasElement} theCanvas
 * @return {(event: MessageEvent) => void}
 */
function handleMessage(theCanvas) {
    return function (event) {
        if (lastRAFid != null) {
            cancelAnimationFrame(lastRAFid)
        }
        assertNotNull(gameSettings)
        const {rows, columns, cellWidth, cellHeight} = gameSettings
        const binData = /** @type {ArrayLike<string>} */atob(event.data)
        const u8arr = Uint8Array.from(binData, b64byte => b64byte.codePointAt(0))
        const padding = u8arr[0]
        const data = u8arr.slice(1)


        console.log(
            Array.from(data, x => x.toString(2).padStart(8, '0'))
                .join('')
                .match(new RegExp(`.{1,${gameSettings.columns}}`, 'g'))
                .map(x => x.match(/.{1,8}/g).join(' '))
                .join('\n')
        )

        lastRAFid = requestAnimationFrame(function () {
            for (let row = 0; row < rows; row++) {
                for (let column = 0; column < columns; column++) {
                    const ctx = theCanvas.getContext("2d")
                    assertNotNull(ctx)

                    ctx.fillStyle = "white";
                    if (isAlive(data, {row, column})) {
                        ctx.fillRect(column * cellWidth, row * cellHeight, cellWidth, cellHeight)
                    } else {
                        ctx.clearRect(column * cellWidth, row * cellHeight, cellWidth, cellHeight)
                    }
                }
            }
        })
    }
}

/**
 * @param {Uint8Array} u8arr
 * @param {number} row
 * @param {number} column
 */
function isAlive(u8arr, {row, column}) {
    assertNotNull(gameSettings)
    const bitPos = row * gameSettings.columns + column
    const bytePos = Math.floor(bitPos / 8)
    const inBytePos = bitPos % 8
    const shiftBy = 8 - inBytePos - 1
    const mask = 1 << shiftBy
    const byteAtPos = u8arr[bytePos]
    return (byteAtPos & mask) !== 0
}

/**
 *
 * @param {unknown} data
 * @return {GameSettings}
 */
function parseGameSettings(data) {
    assertIsString(data)
    const parsed = JSON.parse(data)
    assertIsObject(parsed)
    const {rows, columns, interval} = parsed
    assertPositiveInteger(rows)
    assertPositiveInteger(columns)
    assertPositiveInteger(interval)

    const {clientHeight, clientWidth} = document.documentElement
    const rowsToColumns = rows / columns
    const columnsToRows = columns / rows

    /** @type {number} */
    let cellHeight
    /** @type {number} */
    let cellWidth

    if (clientWidth * rowsToColumns < clientHeight) {
        const height = clientWidth * rowsToColumns
        cellHeight = height / rows
        cellWidth = clientWidth / columns
    } else {
        const width = clientHeight * columnsToRows
        cellWidth = width / columns
        cellHeight = clientHeight / rows
    }

    cellWidth = Math.floor(cellWidth)
    cellHeight = Math.floor(cellHeight)
    const boardHeight = cellHeight * rows
    const boardWidth = cellWidth * columns

    const val = {rows, columns, interval, cellWidth, cellHeight, boardWidth, boardHeight}
    console.log("Game Settings", val, {rowsToColumns, columnsToRows, clientHeight, clientWidth})

    return val
}

class AssertionError extends Error {
}

/**
 *
 * @param {unknown} val
 * @param {string} [msg]
 * @return {asserts val is string}
 */
function assertIsString(val, msg) {
    if (typeof val !== "string") {
        throw new AssertionError(msg ?? `val=${String(val)} is not a string`)
    }
}

/**
 *
 * @param {unknown} val
 * @param {string} [msg]
 * @return {asserts val is Record<string, unknown>}
 */
function assertIsObject(val, msg) {
    assertNotNull(val)
    if (typeof val !== "object") {
        throw new AssertionError(msg ?? `val=${String(val)} is not an object`)
    }
}

/**
 *
 * @template {any} T
 * @param {T} val
 * @param {string} [msg]
 * @return {asserts val is NonNullable<T>}
 */
function assertNotNull(val, msg) {
    if (val == null) {
        throw new AssertionError(msg)
    }
}

/**
 *
 * @template {any} T
 * @param {T} val
 * @param {string} [msg]
 * @return {asserts val is number}
 */
function assertNumber(val, msg) {
    if (typeof val != "number") {
        throw new AssertionError(msg ?? `val="${String(val)}" is not a number`)
    }
}

/**
 * @param {unknown} val
 * @param {string} [msg]
 * @return {asserts val is number}
 */
function assertPositiveInteger(val, msg) {
    assertNumber(val, msg)
    if (!Number.isInteger(val) || val <= 0) {
        throw new AssertionError(msg ?? `val="${String(val)}" is not a positive integer`)
    }
}