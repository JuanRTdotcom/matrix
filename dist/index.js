"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const config_1 = require("./config/config");
const server_1 = require("./server/server");
const cfg = (0, config_1.loadConfig)();
const app = (0, server_1.buildApp)(cfg);
app.listen(Number(cfg.port), () => {
    console.log(`API en Node.js escuchando en :${cfg.port}`);
});
exports.default = app;
//# sourceMappingURL=index.js.map