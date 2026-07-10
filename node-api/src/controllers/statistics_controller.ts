import type { Request, Response } from 'express';
import { ValidationError } from '../usecases/calculate_statistics';
import type { StatisticsService } from './interfaces';
import { successOK, errorResponse } from '../config/response';

/** Controlador HTTP del endpoint de estadísticas. */
export class StatisticsController {
    constructor(private service: StatisticsService) {
        this.calcular = this.calcular.bind(this);
    }

    /**
     * POST /api/v1/statistics
     * Body: { "matrices": number[][][] }
     * Devuelve max, min, promedio, suma total y si alguna matriz es diagonal.
     */
    calcular(req: Request, res: Response): void {
        try {
            const { matrices } = req.body;
            const stats = this.service.execute(matrices);
            res.status(200).json(successOK('Estadísticas calculadas correctamente', stats));
        } catch (err) {
            if (err instanceof ValidationError) {
                res.status(400).json(errorResponse(400, 'bad_request', err.message));
                return;
            }
            res.status(500).json(errorResponse(500, 'internal_error', 'Error al calcular estadísticas'));
        }
    }
}
