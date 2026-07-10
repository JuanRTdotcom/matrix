import type { Matrix, MatrixSet } from '../models/matrix';
import { Statistics } from '../models/statistics';
import type { StatisticsRepository } from '../usecases/calculate_statistics';

/** MatrixRepository calcula las estadísticas sobre un conjunto de matrices. */
export class MatrixRepository implements StatisticsRepository {
    /** Calcula max, min, promedio, suma total y si alguna matriz es diagonal. */
    calcularEstadisticas(matrices: MatrixSet): Statistics {
        const valores: number[] = [];

        for (const matriz of matrices) {
            for (const fila of matriz) {
                for (const valor of fila) {
                    valores.push(valor);
                }
            }
        }

        if (valores.length === 0) {
            throw new Error('El conjunto de matrices no contiene valores');
        }

        const sum = valores.reduce((acc, v) => acc + v, 0);
        const max = Math.max(...valores);
        const min = Math.min(...valores);
        const average = sum / valores.length;
        const isDiagonal = matrices.some((m) => this.esDiagonal(m));

        return new Statistics(max, min, this.redondear(average), this.redondear(sum), isDiagonal);
    }

    /** Una matriz es diagonal si es cuadrada y todo lo que está fuera de la diagonal es cero. */
    private esDiagonal(matriz: Matrix): boolean {
        const filas = matriz.length;
        if (filas === 0) return false;
        const columnas = matriz[0].length;
        if (filas !== columnas) return false;

        const eps = 1e-9;
        for (let i = 0; i < filas; i++) {
            if (matriz[i].length !== columnas) return false;
            for (let j = 0; j < columnas; j++) {
                if (i !== j && Math.abs(matriz[i][j]) > eps) {
                    return false;
                }
            }
        }
        return true;
    }

    /** Redondea a 6 decimales para evitar ruido de punto flotante. */
    private redondear(x: number): number {
        return Math.round(x * 1e6) / 1e6;
    }
}
