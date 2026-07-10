import dotenv from 'dotenv';

dotenv.config();

/** Configuración del servicio leída del entorno. */
export interface Config {
    port: string;
    jwtSecret: string;
    corsOrigins: string;
}

/** Carga la config desde el entorno; si falta una variable requerida, aborta el arranque. */
export function loadConfig(): Config {
    const missing: string[] = [];

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
function mustEnv(key: string, missing: string[]): string {
    const value = process.env[key];
    if (!value) {
        missing.push(key);
        return '';
    }
    return value;
}

/** Lee una variable opcional con valor por defecto. */
function getEnv(key: string, fallback: string): string {
    return process.env[key] ?? fallback;
}
