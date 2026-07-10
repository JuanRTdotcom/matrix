import { CalculateStatisticsUseCase, ValidationError } from './calculate_statistics';
import { MatrixRepository } from '../repositories/matrix_repository';

function buildUseCase(): CalculateStatisticsUseCase {
    return new CalculateStatisticsUseCase(new MatrixRepository());
}

describe('CalculateStatisticsUseCase', () => {
    it('calcula max, min, suma y promedio sobre varias matrices', () => {
        const uc = buildUseCase();
        const stats = uc.execute([
            [
                [1, 2],
                [3, 4]
            ],
            [[5, 6]]
        ]);

        expect(stats.max).toBe(6);
        expect(stats.min).toBe(1);
        expect(stats.sum).toBe(21);
        expect(stats.average).toBe(3.5);
    });

    it('detecta una matriz diagonal', () => {
        const uc = buildUseCase();
        const stats = uc.execute([
            [
                [5, 0],
                [0, 9]
            ]
        ]);
        expect(stats.isDiagonal).toBe(true);
    });

    it('reporta que no es diagonal cuando hay valores fuera de la diagonal', () => {
        const uc = buildUseCase();
        const stats = uc.execute([
            [
                [5, 1],
                [0, 9]
            ]
        ]);
        expect(stats.isDiagonal).toBe(false);
    });

    it('lanza ValidationError con entrada vacía', () => {
        const uc = buildUseCase();
        expect(() => uc.execute([])).toThrow(ValidationError);
    });

    it('lanza ValidationError con valores no numéricos', () => {
        const uc = buildUseCase();
        // @ts-expect-error prueba de entrada inválida
        expect(() => uc.execute([[[1, 'x']]])).toThrow(ValidationError);
    });
});
