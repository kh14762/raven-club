import {PUBLIC_API_URL} from "$env/static/public"

export async function registerUser(data: {
    username: string;
    email: string;
    password: string;
    confirmPassword: string;
}) {
    return await fetch(PUBLIC_API_URL + '/api/auth/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    });
}

export async function loginUser(data: {
    username: string;
    password: string;
}) {
    return await fetch(PUBLIC_API_URL + '/api/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    });
}