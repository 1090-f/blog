import { spawn } from 'node:child_process'
import { fileURLToPath } from 'node:url'

const viteCli = fileURLToPath(new URL('../node_modules/vite/bin/vite.js', import.meta.url))
const processes = [
  spawn(process.execPath, [viteCli], { stdio: 'inherit' }),
  spawn(process.execPath, [viteCli, '--mode', 'admin'], { stdio: 'inherit' })
]

let shuttingDown = false

function shutdown(code = 0) {
  if (shuttingDown) {
    return
  }

  shuttingDown = true
  for (const child of processes) {
    if (!child.killed) {
      child.kill()
    }
  }
  process.exit(code)
}

for (const child of processes) {
  child.on('error', () => shutdown(1))
  child.on('exit', code => {
    if (!shuttingDown && code !== 0) {
      shutdown(code ?? 1)
    }
  })
}

process.on('SIGINT', () => shutdown())
process.on('SIGTERM', () => shutdown())
