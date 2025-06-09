export interface User {
    id: string;
    username: string;
    email: string;
    password_hash: string;
    created_at: Date;
    updated_at: Date;
    is_deleted: boolean;
    deleted_at: Date | null;
}
