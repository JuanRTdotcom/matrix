"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.buildApp = buildApp;
const express_1 = __importDefault(require("express"));
const cors_1 = __importDefault(require("cors"));
const matrix_repository_1 = require("../repositories/matrix_repository");
const calculate_statistics_1 = require("../usecases/calculate_statistics");
const statistics_controller_1 = require("../controllers/statistics_controller");
const routes_1 = require("./routes");
/**
 * buildApp arma la app Express (dependencias, middlewares y rutas) sin ponerla a escuchar.
 * Separar la construcción del `listen` permite testear con supertest sin abrir un puerto.
 */
function buildApp(cfg) {
    // Inyección de dependencias: repository -> usecase -> controller.
    const repository = new matrix_repository_1.MatrixRepository();
    const usecase = new calculate_statistics_1.CalculateStatisticsUseCase(repository);
    const controller = new statistics_controller_1.StatisticsController(usecase);
    // App Express + middlewares globales.
    const app = (0, express_1.default)();
    app.use(express_1.default.json({ limit: '5mb' }));
    app.use((0, cors_1.default)({
        origin: cfg.corsOrigins === '*' ? true : cfg.corsOrigins.split(','),
        methods: ['GET', 'POST']
    }));
    (0, routes_1.registerRoutes)(app, cfg, controller);
    return app;
}
//# sourceMappingURL=server.js.map