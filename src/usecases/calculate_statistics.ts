import type { MatrixSet } from '../models/matrix';
import { Statistics } from '../models/statistics';

/** Dependencia del caso de uso: la fuente que calcula las estadísticas. */
export interface StatisticsRepository {
    calcularEstadisticas(matrices: MatrixSet): Statistics;
}

/** Error de validación de dominio; el controlador lo mapea a un 400. */
export class ValidationError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'ValidationError';
    }
}

/** Caso de uso: calcular estadísticas de un conjunto de matrices. */
export class CalculateStatisticsUseCase {
    constructor(private repository: StatisticsRepository) {}

    execute(matrices: MatrixSet): Statistics {
        this.validar(matrices);
        return this.repository.calcularEstadisticas(matrices);
    }

    /** Reglas de negocio de entrada: al menos una matriz no vacía con números válidos. */
    private validar(matrices: MatrixSet): void {
        if (!Array.isArray(matrices) || matrices.length === 0) {
            throw new ValidationError('El campo "matrices" debe ser un arreglo con al menos una matriz');
        }
        for (const matriz of matrices) {
            if (!Array.isArray(matriz) || matriz.length === 0) {
                throw new ValidationError('Cada matriz debe ser un arreglo de filas no vacío');
            }
            for (const fila of matriz) {
                if (!Array.isArray(fila) || fila.some((v) => typeof v !== 'number' || Number.isNaN(v))) {
                    throw new ValidationError('Todas las filas deben contener únicamente números válidos');
                }
            }
        }
    }
}
