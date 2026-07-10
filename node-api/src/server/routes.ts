import { type Express, Router } from 'express';
import type { Config } from '../config/config';
import type { StatisticsController } from '../controllers/statistics_controller';
import { jwtAuth } from '../middlewares/auth';
import { notFound, errorHandler } from '../middlewares/error_handler';

/** Define todas las rutas en un solo lugar. */
export function registerRoutes(app: Express, cfg: Config, statistics: StatisticsController): void {
    // Público.
    app.get('/', (_req, res) => {
        res.json({
            status: true,
            code: 200,
            message: 'API en Node.js (Express) — Estadísticas de matrices',
            version: '1.0.0'
        });
    });

    // Protegido (requiere JWT).
    const api = Router();
    api.post('/statistics', jwtAuth(cfg.jwtSecret), statistics.calcular); // Estadísticas de un set de matrices
    app.use('/api/v1', api);

    app.use(notFound); // 404 para rutas no registradas
    app.use(errorHandler); // manejo global de errores
}
