"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.StatisticsController = void 0;
const calculate_statistics_1 = require("../usecases/calculate_statistics");
const response_1 = require("../config/response");
/** Controlador HTTP del endpoint de estadísticas. */
class StatisticsController {
    constructor(service) {
        this.service = service;
        this.calcular = this.calcular.bind(this);
    }
    /**
     * POST /api/v1/statistics
     * Body: { "matrices": number[][][] }
     * Devuelve max, min, promedio, suma total y si alguna matriz es diagonal.
     */
    calcular(req, res) {
        try {
            const { matrices } = req.body;
            const stats = this.service.execute(matrices);
            res.status(200).json((0, response_1.successOK)('Estadísticas calculadas correctamente', stats));
        }
        catch (err) {
            if (err instanceof calculate_statistics_1.ValidationError) {
                res.status(400).json((0, response_1.errorResponse)(400, 'bad_request', err.message));
                return;
            }
            res.status(500).json((0, response_1.errorResponse)(500, 'internal_error', 'Error al calcular estadísticas'));
        }
    }
}
exports.StatisticsController = StatisticsController;
//# sourceMappingURL=statistics_controller.js.map