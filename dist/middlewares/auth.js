"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.jwtAuth = jwtAuth;
const jsonwebtoken_1 = __importDefault(require("jsonwebtoken"));
const response_1 = require("../config/response");
/** jwtAuth recibe el secreto y devuelve el middleware que valida el JWT Bearer. */
function jwtAuth(secret) {
    return (req, res, next) => {
        const header = req.header('Authorization');
        if (!header) {
            res.status(401).json((0, response_1.errorResponse)(401, 'unauthorized', 'Header Authorization no encontrado'));
            return;
        }
        if (!header.startsWith('Bearer ')) {
            res.status(401).json((0, response_1.errorResponse)(401, 'unauthorized', 'El header Authorization debe usar el esquema Bearer'));
            return;
        }
        const token = header.substring('Bearer '.length);
        try {
            const payload = jsonwebtoken_1.default.verify(token, secret);
            req.user = typeof payload.sub === 'string' ? payload.sub : undefined;
            next();
        }
        catch {
            res.status(401).json((0, response_1.errorResponse)(401, 'unauthorized', 'Token inválido o expirado'));
        }
    };
}
//# sourceMappingURL=auth.js.map