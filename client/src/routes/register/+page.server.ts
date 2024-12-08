import type {Actions} from './$types'
import {fail} from '@sveltejs/kit'
import {z, ZodError} from "zod";
import { PUBLIC_API_URL } from "$env/static/public"
import { setCookie } from "$lib/cookie";

const passwordSchema = z
    .string()
    .min(8, { message: "Password must be at least 8 characters long" })
    .max(20, { message: "Password too long" })
    .refine((password) => /[A-Z]/.test(password), {
        message: "Password must contain at least 1 uppercase",
    })
    .refine((password) => /[a-z]/.test(password), {
        message: "Password must contain at least 1 lowercase",
    })
    .refine((password) => /[0-9]/.test(password),
        { message: "Password must contain a number" })
    .refine((password) => /[!@#$%^&*]/.test(password), {
        message: "Password must contain a special character (!@#$%^&*)",
    });

const userSchema = z.object({
    username: z.string().min(3, 'username must be at least 3 characters long'),
    email: z.string().email('Invalid email address'),
    password: passwordSchema,
    confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
    message: "Passwords must match",
    path: ["confirmPassword"],
})

export const actions: Actions = {
    register: async ({request}) => {
        const data = await request.formData();
        console.log(data)
        let formEntries = Object.fromEntries(data)
        try {
            const {username, email, password, confirmPassword} = userSchema.parse(formEntries)

            //send data to server
            const response = await fetch(PUBLIC_API_URL + '/api/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    email: email,
                    password: password,
                    confirmPassword: confirmPassword // Send just in-case TODO: hash this
                })
            })

            if (response.ok) {
                const data = await response.json();
                console.log(data)
                // TODO store JWTs in cookies
                setCookie("accessToken", data.accessToken, 1)
                setCookie("refreshToken", data.refreshToken, 7 * 24) // 7 days for refresh token

                // TODO call /api/auth/login to auto login the user

                // TODO route to dashboard page
            } else {
                const error = await response.json();
                console.error('Error registering user:', error);
            }

            // return any errors
        } catch (error: unknown) {
            if (error instanceof ZodError) {
                const errors = error.flatten()
                const { username, email, password, confirmPassword } = formEntries
                const { fieldErrors } = errors
                return fail(400, {
                    username: typeof username === 'string' ? username : '',
                    email: typeof email === 'string' ? email: '',
                    password: typeof password === 'string' ? password : '', // TODO: this should be hashed
                    confirmPassword: typeof confirmPassword === 'string' ? confirmPassword : '', // TODO: this should be hashed
                    error: {
                        ...(fieldErrors?.username ? { field : 'username', message: fieldErrors.username[0]} : {}),
                        ...(fieldErrors?.email ? { field : 'email', message: fieldErrors.email[0]} : {}),
                        ...(fieldErrors?.password ? { field : 'password', message: fieldErrors.password[0]} : {}),
                        ...(fieldErrors?.confirmPassword ? { field : 'confirm', message: fieldErrors.confirmPassword[0]} : {})
                    }
                })
            } else {
                console.error('Error making the request:', error) // TODO: handle error, on +page.svelte
            }
            // rethrow to enable redirect
            throw error // TODO: may not have to do this
        }
    }
} satisfies Actions
