import cluster from "cluster";
import os from "os";
import { createServer } from "http";
import entry from "../dist/entry.server.mjs";

if (cluster.isPrimary) {
  os.cpus().forEach(() => {
    cluster.fork();
  });
} else {
  const server = createServer((req, res) => {
    res.end(entry());
  });
  server.listen(3000);
  console.log("listening on :3000");
}
