"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Statistics = void 0;
/** Estadísticas calculadas sobre un conjunto de matrices. */
class Statistics {
    constructor(max, min, average, sum, isDiagonal) {
        this.max = max;
        this.min = min;
        this.average = average;
        this.sum = sum;
        this.isDiagonal = isDiagonal;
    }
    static empty() {
        return new Statistics(0, 0, 0, 0, false);
    }
}
exports.Statistics = Statistics;
//# sourceMappingURL=statistics.js.map