import express, { type Express } from 'express';
import cors from 'cors';
import type { Config } from '../config/config';
import { MatrixRepository } from '../repositories/matrix_repository';
import { CalculateStatisticsUseCase } from '../usecases/calculate_statistics';
import { StatisticsController } from '../controllers/statistics_controller';
import { registerRoutes } from './routes';

/**
 * buildApp arma la app Express (dependencias, middlewares y rutas) sin ponerla a escuchar.
 * Separar la construcción del `listen` permite testear con supertest sin abrir un puerto.
 */
export function buildApp(cfg: Config): Express {
    // Inyección de dependencias: repository -> usecase -> controller.
    const repository = new MatrixRepository();
    const usecase = new CalculateStatisticsUseCase(repository);
    const controller = new StatisticsController(usecase);

    // App Express + middlewares globales.
    const app = express();
    app.use(express.json({ limit: '5mb' }));
    app.use(
        cors({
            origin: cfg.corsOrigins === '*' ? true : cfg.corsOrigins.split(','),
            methods: ['GET', 'POST']
        })
    );

    registerRoutes(app, cfg, controller);

    return app;
}
