import type { MatrixAdapterInterface } from '$lib/adapters/interfaces/matrix_adapter_interface';
import type { Session } from '$lib/domain/entities/session';

/** Caso de uso: iniciar sesión y obtener el token JWT. */
export class LoginUseCase {
	constructor(private adapter: MatrixAdapterInterface) {}

	async execute(username: string, password: string): Promise<Session> {
		if (!username || !password) {
			throw new Error('Usuario y contraseña son obligatorios');
		}
		return this.adapter.autenticar(username, password);
	}
}
