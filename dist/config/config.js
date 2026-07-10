"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.loadConfig = loadConfig;
const dotenv_1 = __importDefault(require("dotenv"));
dotenv_1.default.config();
/** Carga la config desde el entorno; si falta una variable requerida, aborta el arranque. */
function loadConfig() {
    const missing = [];
    // Requerida.
    const jwtSecret = mustEnv('JWT_SECRET', missing);
    if (missing.length > 0) {
        console.error(`Faltan variables de entorno requeridas: ${missing.join(', ')}`);
        process.exit(1);
    }
    return {
        jwtSecret,
        // Opcionales, con valor por defecto.
        port: getEnv('PORT', '4000'),
        corsOrigins: getEnv('CORS_ORIGINS', '*')
    };
}
/** Lee una variable requerida; si falta, la registra en `missing`. */
function mustEnv(key, missing) {
    const value = process.env[key];
    if (!value) {
        missing.push(key);
        return '';
    }
    return value;
}
/** Lee una variable opcional con valor por defecto. */
function getEnv(key, fallback) {
    return process.env[key] ?? fallback;
}
//# sourceMappingURL=config.js.map