import { env } from '$env/dynamic/public';

/** URL base de la API de Go que consume el navegador. */
export const GO_API_URL = env.PUBLIC_GO_API_URL ?? 'http://localhost:3000';

export const ENDPOINT_LOGIN = `${GO_API_URL}/api/v1/auth/login`;
export const ENDPOINT_PROCESS = `${GO_API_URL}/api/v1/matrix/process`;
