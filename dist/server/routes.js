"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.registerRoutes = registerRoutes;
const express_1 = require("express");
const auth_1 = require("../middlewares/auth");
const error_handler_1 = require("../middlewares/error_handler");
/** Define todas las rutas en un solo lugar. */
function registerRoutes(app, cfg, statistics) {
    // Público.
    app.get('/', (_req, res) => {
        res.json({
            status: true,
            code: 200,
            message: 'API en Node.js (Express) — Estadísticas de matrices',
            version: '1.0.0'
        });
    });
    // Protegido (requiere JWT).
    const api = (0, express_1.Router)();
    api.post('/statistics', (0, auth_1.jwtAuth)(cfg.jwtSecret), statistics.calcular); // Estadísticas de un set de matrices
    app.use('/api/v1', api);
    app.use(error_handler_1.notFound); // 404 para rutas no registradas
    app.use(error_handler_1.errorHandler); // manejo global de errores
}
//# sourceMappingURL=routes.js.map