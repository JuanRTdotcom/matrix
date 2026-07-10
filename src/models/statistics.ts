/** Estadísticas calculadas sobre un conjunto de matrices. */
export class Statistics {
    constructor(
        public max: number,
        public min: number,
        public average: number,
        public sum: number,
        public isDiagonal: boolean
    ) {}

    static empty(): Statistics {
        return new Statistics(0, 0, 0, 0, false);
    }
}
