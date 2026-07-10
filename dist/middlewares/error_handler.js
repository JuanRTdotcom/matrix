"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.notFound = notFound;
exports.errorHandler = errorHandler;
const response_1 = require("../config/response");
/** Middleware 404 para rutas no registradas. */
function notFound(_req, res) {
    res.status(404).json((0, response_1.errorResponse)(404, 'not_found', 'El recurso solicitado no fue encontrado'));
}
/** Middleware global de manejo de errores no controlados. */
function errorHandler(err, _req, res, _next) {
    console.error('[error]', err.message);
    res.status(500).json((0, response_1.errorResponse)(500, 'internal_error', 'Error interno del servidor'));
}
//# sourceMappingURL=error_handler.js.map