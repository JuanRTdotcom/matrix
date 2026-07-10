"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.successOK = successOK;
exports.errorResponse = errorResponse;
function successOK(message, data) {
    return { code: 200, status: 'success', message, data };
}
function errorResponse(code, error, message) {
    return { code, status: false, error, message };
}
//# sourceMappingURL=response.js.map