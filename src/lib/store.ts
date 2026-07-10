import { writable } from 'svelte/store';
import { Session } from '$lib/domain/entities/session';

/** Store reactivo con la sesión actual (token JWT). Persiste en sessionStorage. */
function createSessionStore() {
	const initial = Session.empty();
	const { subscribe, set } = writable<Session>(initial);

	return {
		subscribe,
		login(session: Session) {
			if (typeof sessionStorage !== 'undefined') {
				sessionStorage.setItem('token', session.token);
			}
			set(session);
		},
		logout() {
			if (typeof sessionStorage !== 'undefined') {
				sessionStorage.removeItem('token');
			}
			set(Session.empty());
		},
		restore() {
			if (typeof sessionStorage !== 'undefined') {
				const token = sessionStorage.getItem('token');
				if (token) set(new Session(token, 0));
			}
		}
	};
}

export const sessionStore = createSessionStore();
