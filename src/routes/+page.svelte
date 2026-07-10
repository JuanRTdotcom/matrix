<script lang="ts">
	import { onMount } from 'svelte';
	import { sessionStore } from '$lib/store';
	import { MatrixAdapter } from '$lib/adapters/implements/matrix_adapter';
	import { LoginUseCase } from '$lib/usecases/login';
	import { ProcessMatrixUseCase } from '$lib/usecases/process_matrix';
	import { ProcessResult } from '$lib/domain/entities/process_result';

	// Inyección de dependencias.
	const adapter = new MatrixAdapter();
	const loginUseCase = new LoginUseCase(adapter);
	const processUseCase = new ProcessMatrixUseCase(adapter);

	// Estado de autenticación.
	let token = '';
	let username = 'admin';
	let password = 'admin123';

	// Estado de la matriz.
	let matrixText = '12, -51, 4\n6, 167, -68\n-4, 24, -41';
	let result: ProcessResult | null = null;

	let loading = false;
	let error = '';

	onMount(() => {
		sessionStore.restore();
	});

	sessionStore.subscribe((s) => (token = s.token));

	function parseMatrix(text: string): number[][] {
		return text
			.trim()
			.split('\n')
			.filter((line) => line.trim().length > 0)
			.map((line) =>
				line
					.split(/[,\s]+/)
					.filter((n) => n.length > 0)
					.map((n) => {
						const value = Number(n);
						if (Number.isNaN(value)) throw new Error(`Valor no numérico: "${n}"`);
						return value;
					})
			);
	}

	async function handleLogin() {
		error = '';
		loading = true;
		try {
			const session = await loginUseCase.execute(username, password);
			sessionStore.login(session);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Error de autenticación';
		} finally {
			loading = false;
		}
	}

	async function handleProcess() {
		error = '';
		result = null;
		loading = true;
		try {
			const matrix = parseMatrix(matrixText);
			result = await processUseCase.execute(token, matrix);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Error al procesar la matriz';
		} finally {
			loading = false;
		}
	}

	function logout() {
		sessionStore.logout();
		result = null;
	}

	function fmt(n: number): string {
		return Number.isInteger(n) ? String(n) : n.toFixed(4);
	}
</script>

<main>
	<header>
		<h1>Reto Técnico · Interseguro</h1>
		<p>Factorización QR (Go · Fiber) → Estadísticas (Node · Express)</p>
	</header>

	{#if !token}
		<section class="card">
			<h2>Iniciar sesión</h2>
			<label>
				Usuario
				<input bind:value={username} placeholder="admin" />
			</label>
			<label>
				Contraseña
				<input type="password" bind:value={password} placeholder="admin123" />
			</label>
			<button on:click={handleLogin} disabled={loading}>
				{loading ? 'Autenticando…' : 'Entrar'}
			</button>
			<small>Credenciales de demo: admin / admin123</small>
		</section>
	{:else}
		<section class="card">
			<div class="row-between">
				<h2>Matriz de entrada</h2>
				<button class="ghost" on:click={logout}>Cerrar sesión</button>
			</div>
			<p class="hint">Una fila por línea. Separa los valores con comas o espacios.</p>
			<textarea bind:value={matrixText} rows="6"></textarea>
			<button on:click={handleProcess} disabled={loading}>
				{loading ? 'Procesando…' : 'Procesar matriz'}
			</button>
		</section>
	{/if}

	{#if error}
		<p class="error">⚠ {error}</p>
	{/if}

	{#if result}
		<section class="card">
			<h2>Resultados</h2>

			<div class="grid">
				<div>
					<h3>Matriz original</h3>
					<div class="matrix">
						{#each result.original as row}
							<div class="mrow">
								{#each row as val}<span>{fmt(val)}</span>{/each}
							</div>
						{/each}
					</div>
				</div>

				<div>
					<h3>Q (ortonormal)</h3>
					<div class="matrix">
						{#each result.qr.q as row}
							<div class="mrow">
								{#each row as val}<span>{fmt(val)}</span>{/each}
							</div>
						{/each}
					</div>
				</div>

				<div>
					<h3>R (triangular superior)</h3>
					<div class="matrix">
						{#each result.qr.r as row}
							<div class="mrow">
								{#each row as val}<span>{fmt(val)}</span>{/each}
							</div>
						{/each}
					</div>
				</div>
			</div>

			<h3>Estadísticas (calculadas en Node.js sobre Q y R)</h3>
			<div class="stats">
				<div class="stat"><span>Máximo</span><strong>{fmt(result.statistics.max)}</strong></div>
				<div class="stat"><span>Mínimo</span><strong>{fmt(result.statistics.min)}</strong></div>
				<div class="stat"><span>Promedio</span><strong>{fmt(result.statistics.average)}</strong></div>
				<div class="stat"><span>Suma total</span><strong>{fmt(result.statistics.sum)}</strong></div>
				<div class="stat">
					<span>¿Alguna diagonal?</span>
					<strong class:yes={result.statistics.isDiagonal}>
						{result.statistics.isDiagonal ? 'Sí' : 'No'}
					</strong>
				</div>
			</div>
		</section>
	{/if}

	<footer>API en Go :3000 · API en Node :4000 · Comunicación HTTP autenticada con JWT</footer>
</main>

<style>
	main {
		max-width: 900px;
		margin: 0 auto;
		padding: 2rem 1.25rem 4rem;
	}
	header h1 {
		margin: 0 0 0.25rem;
		font-size: 1.8rem;
	}
	header p {
		margin: 0 0 1.5rem;
		color: var(--muted);
	}
	.card {
		background: var(--panel);
		border: 1px solid var(--border);
		border-radius: 12px;
		padding: 1.5rem;
		margin-bottom: 1.25rem;
	}
	h2 {
		margin: 0 0 1rem;
		font-size: 1.2rem;
	}
	h3 {
		font-size: 0.95rem;
		color: var(--accent-2);
		margin: 1rem 0 0.5rem;
	}
	label {
		display: block;
		margin-bottom: 0.85rem;
		font-size: 0.9rem;
		color: var(--muted);
	}
	input,
	textarea {
		width: 100%;
		margin-top: 0.35rem;
		padding: 0.6rem 0.75rem;
		background: var(--panel-2);
		border: 1px solid var(--border);
		border-radius: 8px;
		color: var(--text);
		font-size: 0.95rem;
		font-family: 'Consolas', 'Menlo', monospace;
	}
	button {
		background: var(--accent);
		color: #06240f;
		border: none;
		border-radius: 8px;
		padding: 0.65rem 1.25rem;
		font-weight: 600;
		cursor: pointer;
		font-size: 0.95rem;
	}
	button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}
	button.ghost {
		background: transparent;
		color: var(--muted);
		border: 1px solid var(--border);
	}
	small,
	.hint {
		display: block;
		color: var(--muted);
		font-size: 0.8rem;
		margin-top: 0.75rem;
	}
	.hint {
		margin: 0 0 0.75rem;
	}
	.row-between {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	.error {
		color: var(--danger);
		background: rgba(239, 68, 68, 0.1);
		border: 1px solid var(--danger);
		border-radius: 8px;
		padding: 0.75rem 1rem;
	}
	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
		gap: 1rem;
	}
	.matrix {
		background: var(--panel-2);
		border-radius: 8px;
		padding: 0.6rem;
		font-family: 'Consolas', 'Menlo', monospace;
		font-size: 0.85rem;
		overflow-x: auto;
	}
	.mrow {
		display: flex;
		gap: 0.5rem;
	}
	.mrow span {
		min-width: 3.5rem;
		text-align: right;
	}
	.stats {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
		gap: 0.75rem;
		margin-top: 0.5rem;
	}
	.stat {
		background: var(--panel-2);
		border-radius: 8px;
		padding: 0.75rem;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}
	.stat span {
		color: var(--muted);
		font-size: 0.8rem;
	}
	.stat strong {
		font-size: 1.15rem;
	}
	.stat strong.yes {
		color: var(--accent);
	}
	footer {
		color: var(--muted);
		font-size: 0.8rem;
		text-align: center;
		margin-top: 2rem;
	}
</style>
