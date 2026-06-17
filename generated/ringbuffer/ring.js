#!/usr/bin/env bun
// Data-oriented ring buffer in Bun/JS.
// No structs, no classes, no type system nonsense.
// Flat array, scalar head/count, freestanding functions.

const CAPACITY = 1024;

function ringPush(buf, head, count, value) {
    buf[head] = value;
    head = (head + 1) % CAPACITY;
    if (count < CAPACITY) count++;
    return [head, count];
}

function ringAvg(buf, count) {
    if (count === 0) return 0.0;
    let sum = 0.0;
    for (let i = 0; i < count; i++) sum += buf[i];
    return sum / count;
}

function ringReadRange(buf, head, count, n) {
    if (n <= 0 || count === 0) return [];
    if (n > count) n = count;
    const res = [];
    for (let i = 0; i < n; i++) {
        res.push(buf[(head - n + i + CAPACITY) % CAPACITY]);
    }
    return res;
}

function ringUtilization(count) {
    return count / CAPACITY;
}

async function main() {
    const storage = new Array(CAPACITY).fill(0.0);
    let head = 0, count = 0;

    const text = await Bun.stdin.text();
    for (const raw of text.split('\n')) {
        const line = raw.trim();
        if (!line) continue;
        const parts = line.split(/\s+/);
        const cmd = parts[0];
        if (!cmd) continue;

        if (cmd === 'push') {
            if (parts.length < 2) continue;
            const value = parseFloat(parts[1]);
            [head, count] = ringPush(storage, head, count, value);
        } else if (cmd === 'avg') {
            console.log(ringAvg(storage, count).toFixed(1));
        } else if (cmd === 'readRange') {
            const n = parts.length > 1 ? parseInt(parts[1], 10) : 1;
            for (const v of ringReadRange(storage, head, count, n)) {
                console.log(v.toFixed(1));
            }
        } else if (cmd === 'utilization') {
            console.log(ringUtilization(count).toFixed(4));
        }
    }
}

main();
