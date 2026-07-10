import type { ProcessResult } from '$lib/domain/entities/process_result';
import type { Session } from '$lib/domain/entities/session';

/** Contrato del adaptador que traduce datos remotos a entidades de dominio. */
export interface MatrixAdapterInterface {
	autenticar(username: string, password: string): Promise<Session>;
	procesarMatriz(token: string, matrix: number[][]): Promise<ProcessResult>;
}
