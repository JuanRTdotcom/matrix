import request from 'supertest';
import jwt from 'jsonwebtoken';
import { buildApp } from './server/server';
import type { Config } from './config/config';

// La config se construye a mano en el test (igual que en el test de integración de Go),
// sin depender del entorno.
const TEST_SECRET = 'test-secret';
const cfg: Config = { port: '0', jwtSecret: TEST_SECRET, corsOrigins: '*' };
const app = buildApp(cfg);

/** Firma un JWT válido con el mismo secreto que valida el middleware. */
function validToken(): string {
    return jwt.sign({ sub: 'admin' }, TEST_SECRET, { expiresIn: '1h' });
}

describe('POST /api/v1/statistics (integración)', () => {
    it('calcula estadísticas con un JWT válido', async () => {
        const res = await request(app)
            .post('/api/v1/statistics')
            .set('Authorization', `Bearer ${validToken()}`)
            .send({ matrices: [[[1, 2], [3, 4]], [[5, 6]]] });

        expect(res.status).toBe(200);
        expect(res.body.data.max).toBe(6);
        expect(res.body.data.min).toBe(1);
        expect(res.body.data.sum).toBe(21);
        expect(res.body.data.average).toBe(3.5);
    });

    it('detecta una matriz diagonal', async () => {
        const res = await request(app)
            .post('/api/v1/statistics')
            .set('Authorization', `Bearer ${validToken()}`)
            .send({ matrices: [[[5, 0], [0, 9]]] });

        expect(res.status).toBe(200);
        expect(res.body.data.isDiagonal).toBe(true);
    });

    it('rechaza la petición sin token (401)', async () => {
        const res = await request(app)
            .post('/api/v1/statistics')
            .send({ matrices: [[[1]]] });

        expect(res.status).toBe(401);
    });

    it('rechaza un token inválido (401)', async () => {
        const res = await request(app)
            .post('/api/v1/statistics')
            .set('Authorization', 'Bearer token-falso')
            .send({ matrices: [[[1]]] });

        expect(res.status).toBe(401);
    });

    it('rechaza un body inválido con 400', async () => {
        const res = await request(app)
            .post('/api/v1/statistics')
            .set('Authorization', `Bearer ${validToken()}`)
            .send({ matrices: [] });

        expect(res.status).toBe(400);
    });

    it('responde 404 en rutas inexistentes', async () => {
        const res = await request(app).get('/api/v1/nope');
        expect(res.status).toBe(404);
    });
});
