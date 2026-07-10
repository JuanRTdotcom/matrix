import type { Request, Response, NextFunction } from 'express';
import jwt from 'jsonwebtoken';
import { errorResponse } from '../config/response';

/** Extiende Request para exponer el usuario autenticado. */
export interface AuthRequest extends Request {
    user?: string;
}

/** jwtAuth recibe el secreto y devuelve el middleware que valida el JWT Bearer. */
export function jwtAuth(secret: string) {
    return (req: AuthRequest, res: Response, next: NextFunction): void => {
        const header = req.header('Authorization');

        if (!header) {
            res.status(401).json(errorResponse(401, 'unauthorized', 'Header Authorization no encontrado'));
            return;
        }
        if (!header.startsWith('Bearer ')) {
            res.status(401).json(errorResponse(401, 'unauthorized', 'El header Authorization debe usar el esquema Bearer'));
            return;
        }

        const token = header.substring('Bearer '.length);

        try {
            const payload = jwt.verify(token, secret) as jwt.JwtPayload;
            req.user = typeof payload.sub === 'string' ? payload.sub : undefined;
            next();
        } catch {
            res.status(401).json(errorResponse(401, 'unauthorized', 'Token inválido o expirado'));
        }
    };
}
