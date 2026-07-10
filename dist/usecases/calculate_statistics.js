"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CalculateStatisticsUseCase = exports.ValidationError = void 0;
/** Error de validación de dominio; el controlador lo mapea a un 400. */
class ValidationError extends Error {
    constructor(message) {
        super(message);
        this.name = 'ValidationError';
    }
}
exports.ValidationError = ValidationError;
/** Caso de uso: calcular estadísticas de un conjunto de matrices. */
class CalculateStatisticsUseCase {
    constructor(repository) {
        this.repository = repository;
    }
    execute(matrices) {
        this.validar(matrices);
        return this.repository.calcularEstadisticas(matrices);
    }
    /** Reglas de negocio de entrada: al menos una matriz no vacía con números válidos. */
    validar(matrices) {
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
exports.CalculateStatisticsUseCase = CalculateStatisticsUseCase;
//# sourceMappingURL=calculate_statistics.js.map