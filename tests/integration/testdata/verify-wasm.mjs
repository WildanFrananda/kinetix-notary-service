import { readFileSync } from "node:fs";

const [wasmExec, module, document, settings] = process.argv.slice(2);

await import(wasmExec);

const go = new globalThis.Go();
const { instance } = await WebAssembly.instantiate(readFileSync(module), go.importObject);
go.run(instance);

const answer = await globalThis.kinetixNotaryVerify(readFileSync(document, "utf8"), JSON.parse(settings));

process.stdout.write(answer.stdout);
process.stderr.write(answer.stderr);
process.exit(answer.exit);
