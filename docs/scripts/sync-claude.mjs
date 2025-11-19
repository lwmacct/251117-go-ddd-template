import { promises as fs } from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const repoRoot = path.resolve(__dirname, "../..");
const sourcePath = path.join(repoRoot, "CLAUDE.md");
const destinationPath = path.resolve(__dirname, "../development/claude-guide.md");

const banner = [
  "<!--",
  "  ⚠️ 该文件由 docs/scripts/sync-claude.mjs 自动生成。",
  "  请在仓库根目录 CLAUDE.md 中修改内容，然后运行：",
  "    npm --prefix docs run sync:claude",
  "-->",
  "",
].join("\n");

async function main() {
  try {
    const content = await fs.readFile(sourcePath, "utf8");
    await fs.mkdir(path.dirname(destinationPath), { recursive: true });
    await fs.writeFile(destinationPath, `${banner}${content}`, "utf8");
    console.log(
      `Synced ${path.relative(repoRoot, sourcePath)} -> ${path.relative(repoRoot, destinationPath)}`
    );
  } catch (error) {
    console.error("Failed to sync CLAUDE.md:", error);
    process.exitCode = 1;
  }
}

main();
