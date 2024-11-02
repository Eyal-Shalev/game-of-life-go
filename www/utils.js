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
// export let gameSettings = null

/**
 * @param {GameSettings} gameSettings
 * @param {Uint8Array} u8arr
 * @param {number} row
 * @param {number} column
 */
export function isAlive(gameSettings, u8arr, {row, column}) {
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
export function parseGameSettings(data) {
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
export function assertIsString(val, msg) {
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
export function assertIsObject(val, msg) {
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
export function assertNotNull(val, msg) {
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
export function assertNumber(val, msg) {
    if (typeof val != "number") {
        throw new AssertionError(msg ?? `val="${String(val)}" is not a number`)
    }
}

/**
 * @param {unknown} val
 * @param {string} [msg]
 * @return {asserts val is number}
 */
export function assertPositiveInteger(val, msg) {
    assertNumber(val, msg)
    if (!Number.isInteger(val) || val <= 0) {
        throw new AssertionError(msg ?? `val="${String(val)}" is not a positive integer`)
    }
}

/**
 * @param {unknown} raw
 * @returns {number}
 * @throws {Error} if raw is not a valid number string.
 */
export function parseNumber(raw) {
    const num = Number(raw);
    if (Number.isNaN(num)) {
        throw new Error(`rows="${String(raw)}" is not a number`)
    }
    return num
}

/**
 * @param {number} val
 */
export function validatePositiveInteger(val) {
    if (val <= 0 || val % 1 !== 0) {
        throw new Error(`${val} is not a positive integer`)
    }
}