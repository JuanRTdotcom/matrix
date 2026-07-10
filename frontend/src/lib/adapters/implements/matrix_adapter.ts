import { ApiRepository } from '$lib/repository/api_repository';
import { ProcessResult } from '$lib/domain/entities/process_result';
import { Session } from '$lib/domain/entities/session';
import type { MatrixAdapterInterface } from '../interfaces/matrix_adapter_interface';

/** MatrixAdapter mapea las respuestas de la API a entidades de dominio. */
export class MatrixAdapter implements MatrixAdapterInterface {
	private repository: ApiRepository;

	constructor(repository: ApiRepository = new ApiRepository()) {
		this.repository = repository;
	}

	async autenticar(username: string, password: string): Promise<Session> {
		const { token, expiresIn } = await this.repository.login(username, password);
		return new Session(token, expiresIn);
	}

	async procesarMatriz(token: string, matrix: number[][]): Promise<ProcessResult> {
		const data = await this.repository.procesarMatriz(token, matrix);
		return new ProcessResult(data.original, data.qr, data.statistics);
	}
}
