import type { Request, Response, NextFunction } from 'express';
import { errorResponse } from '../config/response';

/** Middleware 404 para rutas no registradas. */
export function notFound(_req: Request, res: Response): void {
    res.status(404).json(errorResponse(404, 'not_found', 'El recurso solicitado no fue encontrado'));
}

/** Middleware global de manejo de errores no controlados. */
export function errorHandler(err: Error, _req: Request, res: Response, _next: NextFunction): void {
    console.error('[error]', err.message);
    res.status(500).json(errorResponse(500, 'internal_error', 'Error interno del servidor'));
}
