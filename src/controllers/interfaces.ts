import type { MatrixSet } from '../models/matrix';
import type { Statistics } from '../models/statistics';

/** Abstracción del caso de uso que usa el controlador, para no depender de su implementación. */
export interface StatisticsService {
    execute(matrices: MatrixSet): Statistics;
}
