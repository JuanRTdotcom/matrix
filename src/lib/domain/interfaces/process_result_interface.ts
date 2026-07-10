/** Estadísticas devueltas por la cadena Go -> Node. */
export interface StatisticsInterface {
	max: number;
	min: number;
	average: number;
	sum: number;
	isDiagonal: boolean;
}

/** Factorización QR devuelta por la API en Go. */
export interface QRInterface {
	q: number[][];
	r: number[][];
}

/** Resultado del procesamiento de una matriz: original, factorización QR y estadísticas. */
export interface ProcessResultInterface {
	original: number[][];
	qr: QRInterface;
	statistics: StatisticsInterface;
}
