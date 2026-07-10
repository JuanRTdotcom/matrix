import type { MatrixAdapterInterface } from '$lib/adapters/interfaces/matrix_adapter_interface';
import type { ProcessResult } from '$lib/domain/entities/process_result';

/** Caso de uso: procesar una matriz (factorización QR + estadísticas). */
export class ProcessMatrixUseCase {
	constructor(private adapter: MatrixAdapterInterface) {}

	async execute(token: string, matrix: number[][]): Promise<ProcessResult> {
		if (!token) {
			throw new Error('Sesión no válida. Inicie sesión nuevamente.');
		}
		if (!matrix.length || !matrix[0]?.length) {
			throw new Error('La matriz no puede estar vacía');
		}
		const cols = matrix[0].length;
		if (matrix.some((row) => row.length !== cols)) {
			throw new Error('La matriz debe ser rectangular (todas las filas del mismo tamaño)');
		}
		return this.adapter.procesarMatriz(token, matrix);
	}
}
