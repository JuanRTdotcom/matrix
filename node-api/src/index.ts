import { loadConfig } from './config/config';
import { buildApp } from './server/server';

const cfg = loadConfig();
const app = buildApp(cfg);

app.listen(Number(cfg.port), () => {
    console.log(`API en Node.js escuchando en :${cfg.port}`);
});

export default app;
