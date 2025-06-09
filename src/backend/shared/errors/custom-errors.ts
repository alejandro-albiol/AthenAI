// src/backend/shared/errors/custom-errors.ts
export class BaseError extends Error {
    constructor(
        public override readonly message: string,
        public readonly code: string,
        public readonly statusCode: number
    ) {
        super(message);
        this.name = this.constructor.name;
    }
}

export class UserError extends BaseError {
    constructor(message: string, code: string) {
        super(message, code, 400);
    }
}

export class NotFoundError extends BaseError {
    constructor(message: string) {
        super(message, 'NOT_FOUND', 404);
    }
}

export class ConflictError extends BaseError {
    constructor(message: string) {
        super(message, 'CONFLICT', 409);
    }
}

export class ValidationError extends BaseError {
    constructor(message: string) {
        super(message, 'VALIDATION_ERROR', 400);
    }
}