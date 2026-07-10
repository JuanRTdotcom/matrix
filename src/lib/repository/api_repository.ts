import { ENDPOINT_LOGIN, ENDPOINT_PROCESS } from '$lib/enviroment/constants';
import type { ProcessResultInterface } from '$lib/domain/interfaces/process_result_interface';

/** ApiRepository encapsula las llamadas HTTP a la API de Go. */
export class ApiRepository {
	/** Autentica y devuelve el token JWT y su expiración. */
	async login(username: string, password: string): Promise<{ token: string; expiresIn: number }> {
		const res = await fetch(ENDPOINT_LOGIN, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username, password })
		});

		const body = await res.json();
		if (!res.ok) {
			throw new Error(body?.message ?? 'Error de autenticación');
		}
		return { token: body.data.token, expiresIn: body.data.expiresIn };
	}

	/** Envía la matriz a la API en Go (que a su vez llama a Node) y devuelve el resultado. */
	async procesarMatriz(token: string, matrix: number[][]): Promise<ProcessResultInterface> {
		const res = await fetch(ENDPOINT_PROCESS, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify({ matrix })
		});

		const body = await res.json();
		if (!res.ok) {
			throw new Error(body?.message ?? 'Error al procesar la matriz');
		}
		return body.data as ProcessResultInterface;
	}
}
