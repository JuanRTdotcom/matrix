/** Sesión del usuario autenticado (token JWT emitido por la API en Go). */
export class Session {
	constructor(
		public token: string,
		public expiresIn: number
	) {}

	static empty(): Session {
		return new Session('', 0);
	}

	get isAuthenticated(): boolean {
		return this.token.length > 0;
	}
}
