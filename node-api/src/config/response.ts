/** Respuesta exitosa estándar. */
export interface ApiSuccess<T> {
    code: number;
    status: string;
    message: string;
    data?: T;
}

/** Respuesta de error estándar. */
export interface ApiError {
    code: number;
    status: boolean;
    error: string;
    message: string;
}

export function successOK<T>(message: string, data: T): ApiSuccess<T> {
    return { code: 200, status: 'success', message, data };
}

export function errorResponse(code: number, error: string, message: string): ApiError {
    return { code, status: false, error, message };
}
