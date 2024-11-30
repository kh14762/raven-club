import type {Actions} from './$types'
import {fail} from '@sveltejs/kit'
import {z, ZodError} from "zod";

const userSchema = z.object({
    username: z.string().min(1, 'Name must be at least 1 characters long'),
    email: z.string().email('Invalid email address'),
    password: z.string().min(1, 'Password must be at least 1 characters long'),
    confirmPassword: z.string().min(1, 'Confirm password must be at least 1 characters long'),
}).refine((data) => data.password !== data.confirmPassword, {
    message: "Passwords must match",
    path: ["confirmPassword"], // path of error
})

export const actions: Actions = {
    register: async ({request}) => {
        const data = await request.formData();
        console.log(data)
        let formEntries = Object.fromEntries(data)

        try {
            const {username, email, password, confirmPassword} = userSchema.parse(formEntries)
            console.log(username, email, password, confirmPassword)

            // send data to server
            // const response = await fetch(PUBLIC_API_URL + '/auth/register', {
            //     method: 'POST',
            //     headers: {
            //         Accept: 'application/json',
            //         'Content-Type': 'application/json',
            //     },
            //     body: JSON.stringify({
            //         username: username,
            //         email: email,
            //         password: password, // TODO: hash this
            //         confirmPassword: confirmPassword // Send just in-case TODO: hash this
            //     })
            // })
            //
            // const json = await response.json();
            // console.log(json) // TODO: take actions based on this response such as navigation or auto login

            // return any errors
        } catch (error: unknown) {
            if (error instanceof ZodError) {
                const errors = error.flatten()
                const { username, email, password, confirmPassword } = formEntries
                const { fieldErrors } = errors
                return fail(400, {
                    username: typeof username === 'string' ? username : '',
                    email: typeof email === 'string' ? email: '',
                    password: typeof password === 'string' ? password : '', // TODO: not sure if i need to return this
                    confirmPassword: typeof confirmPassword === 'string' ? confirmPassword : '', // TODO: not sure if i need to return this
                    error: {
                        ...(fieldErrors?.username ? { field : 'username', message: fieldErrors.username[0]} : {}),
                        ...(fieldErrors?.email ? { field : 'email', message: fieldErrors.email[0]} : {}),
                        ...(fieldErrors?.password ? { field : 'password', message: fieldErrors.password[0]} : {}),
                        ...(fieldErrors?.confirmPassword ? { field : 'confirm', message: fieldErrors.confirmPassword[0]} : {})
                    }
                })
            }
            // rethrow to enable redirect
            throw error // TODO: may not have to do this
        }
    }
} satisfies Actions
