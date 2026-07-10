import type {
	ProcessResultInterface,
	QRInterface,
	StatisticsInterface
} from '../interfaces/process_result_interface';

/** Entidad del resultado de procesamiento de una matriz. */
export class ProcessResult implements ProcessResultInterface {
	constructor(
		public original: number[][],
		public qr: QRInterface,
		public statistics: StatisticsInterface
	) {}

	static empty(): ProcessResult {
		return new ProcessResult(
			[],
			{ q: [], r: [] },
			{ max: 0, min: 0, average: 0, sum: 0, isDiagonal: false }
		);
	}
}
